// frontend/js/auth.js
import { apiRequest, setToken, setUser } from "./api.js";

// ⚠️ GANTI dengan Google Client ID kamu
const GOOGLE_CLIENT_ID = "523452215868-hf97c1mqdcue2pu3j7r8avg7q2tqqg4m.apps.googleusercontent.com";

window.onload = () => {
  // kalau sudah login, langsung masuk dashboard
  const existingToken = localStorage.getItem("token");
  if (existingToken) {
    window.location.href = "main_cards.html";
    return;
  }

  // init Google login
  google.accounts.id.initialize({
    client_id: GOOGLE_CLIENT_ID,
    callback: handleCredentialResponse,
  });

  // render button
  google.accounts.id.renderButton(
    document.getElementById("googleBtn"),
    {
      theme: "outline",
      size: "large",
      text: "continue_with",
      shape: "pill",
      width: 320,
    }
  );

  // optional: tampilkan popup one-tap
  // google.accounts.id.prompt();
};

// callback dari google
async function handleCredentialResponse(response) {
  try {
    console.log("Google ID Token:", response.credential);

    // kirim id_token ke backend
    const result = await apiRequest("/auth/google", "POST", {
      id_token: response.credential,
    });

    // backend harus return { token, user }
    if (!result.token) {
      throw new Error("Backend tidak mengembalikan token!");
    }

    setToken(result.token);
    setUser(result.user);

    alert("Login sukses ✅");
    window.location.href = "main_cards.html";

  } catch (err) {
    console.error(err);
    alert("Login gagal ❌: " + err.message);
  }
}
