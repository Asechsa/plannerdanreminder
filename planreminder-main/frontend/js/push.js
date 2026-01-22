// frontend/js/push.js
import { apiRequest } from "./api.js";

// isi dari ENV kamu
const VAPID_PUBLIC_KEY = "BAsWmoz_BnUuAVWw4JU4gbqI5vAP0mjpwq9A2c5-y47E9OZxdoes4F4q1V8YOlRb3ZJ4iKZv9AG0_yLlIAxtSFo";

// convert base64url -> Uint8Array
function urlBase64ToUint8Array(base64String) {
  const padding = "=".repeat((4 - (base64String.length % 4)) % 4);
  const base64 = (base64String + padding).replace(/-/g, "+").replace(/_/g, "/");
  const rawData = atob(base64);
  const outputArray = new Uint8Array(rawData.length);
  for (let i = 0; i < rawData.length; ++i) outputArray[i] = rawData.charCodeAt(i);
  return outputArray;
}

export async function registerServiceWorker() {
  if (!("serviceWorker" in navigator)) throw new Error("Browser tidak support Service Worker");
  // scope: /frontend/ (kalau live server)
  const reg = await navigator.serviceWorker.register("../sw.js");
  return reg;
}

export async function enablePush() {
  if (!("PushManager" in window)) throw new Error("Browser tidak support Push");

  const permission = await Notification.requestPermission();
  if (permission !== "granted") throw new Error("Permission notifikasi ditolak");

  const reg = await registerServiceWorker();

  // cek subscription existing
  let sub = await reg.pushManager.getSubscription();
  if (!sub) {
    sub = await reg.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey: urlBase64ToUint8Array(VAPID_PUBLIC_KEY),
    });
  }

  // kirim ke backend untuk disimpan
  // format sub sudah ada endpoint, keys, endpoint dll.
  await apiRequest("/push/subscribe", "POST", sub);

  return sub;
}

export async function disablePush() {
  const reg = await registerServiceWorker();
  const sub = await reg.pushManager.getSubscription();
  if (!sub) return;

  // inform backend supaya hapus subscription
  await apiRequest("/push/unsubscribe", "POST", { endpoint: sub.endpoint });

  await sub.unsubscribe();
}
