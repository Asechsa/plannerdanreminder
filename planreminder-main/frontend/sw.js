// frontend/sw.js

self.addEventListener("install", (event) => {
  self.skipWaiting();
});

self.addEventListener("activate", (event) => {
  event.waitUntil(self.clients.claim());
});

// terima push dari server
self.addEventListener("push", (event) => {
  let data = {};
  try {
    data = event.data ? event.data.json() : {};
  } catch (e) {
    data = { title: "PlanReminder", body: event.data?.text() || "Reminder masuk!" };
  }

  const title = data.title || "PlanReminder";
  const options = {
    body: data.body || "Reminder masuk!",
    icon: data.icon || "/assets/logo.png",   // optional
    badge: data.badge || "/assets/logo.png", // optional
    data: {
      url: data.url || "/pages/main_cards.html"
    }
  };

  event.waitUntil(self.registration.showNotification(title, options));
});

// klik notifikasi -> buka halaman
self.addEventListener("notificationclick", (event) => {
  event.notification.close();
  const url = event.notification.data?.url || "/pages/main_cards.html";

  event.waitUntil(
    self.clients.matchAll({ type: "window", includeUncontrolled: true }).then((clientList) => {
      for (const client of clientList) {
        if (client.url.includes(url) && "focus" in client) return client.focus();
      }
      if (self.clients.openWindow) return self.clients.openWindow(url);
    })
  );
});
