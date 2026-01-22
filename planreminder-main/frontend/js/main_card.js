// frontend/js/main_cards.js
import { apiRequest, logout, getUser } from "./api.js";

let cards = [];
let selectedCardId = null;

const presetColors = [
  "#3B82F6", "#8B5CF6", "#10B981", "#F97316",
  "#EF4444", "#14B8A6", "#0EA5E9", "#64748B"
];

document.getElementById("logoutBtn")?.addEventListener("click", logout);

async function loadCards() {
  try {
    cards = await apiRequest("/cards");
    renderCards();
  } catch (err) {
    console.error(err.message);
  }
}

function renderCards() {
  const grid = document.getElementById("cardsGrid");
  grid.innerHTML = "";

  cards.forEach(card => {
    const el = document.createElement("div");
    el.className = "relative rounded-2xl p-5 text-white shadow-md hover:shadow-lg transition cursor-pointer";
    el.style.background = card.color || "#3B82F6";

    el.innerHTML = `
      <div class="flex items-start justify-between">
        <div>
          <h3 class="font-bold text-lg">${card.title}</h3>
          <p class="text-xs opacity-90 mt-2">${card.sub_count || 0} Sub-cards • ${card.task_count || 0} Tasks</p>
        </div>

        <button class="menuBtn w-8 h-8 flex items-center justify-center rounded-full hover:bg-white/20 transition" data-id="${card.id}">
          ⋮
        </button>
      </div>

      <div class="dropdown hidden absolute right-4 top-14 bg-white text-slate-800 rounded-xl shadow-lg w-44 overflow-hidden z-20" id="dropdown-${card.id}">
        <button class="w-full text-left px-4 py-2 hover:bg-slate-100 text-sm" onclick="window.openEdit('${card.id}')">Edit Judul</button>
        <button class="w-full text-left px-4 py-2 hover:bg-slate-100 text-sm" onclick="window.openColor('${card.id}')">Ubah Warna</button>
        <button class="w-full text-left px-4 py-2 hover:bg-red-50 text-red-600 text-sm" onclick="window.deleteCard('${card.id}')">Hapus</button>
      </div>
    `;

    // klik card -> sub_cards.html
    el.addEventListener("click", () => {
      window.location.href = `sub_cards.html?card_id=${card.id}`;
    });

    grid.appendChild(el);
  });

  document.querySelectorAll(".menuBtn").forEach(btn => {
    btn.addEventListener("click", (e) => {
      e.stopPropagation();
      toggleDropdown(btn.dataset.id);
    });
  });

  document.addEventListener("click", () => {
    document.querySelectorAll(".dropdown").forEach(d => d.classList.add("hidden"));
  });
}

function toggleDropdown(id) {
  document.querySelectorAll(".dropdown").forEach(d => d.classList.add("hidden"));
  document.getElementById("dropdown-" + id)?.classList.toggle("hidden");
}

// ================================
// Add Card
// ================================
document.getElementById("btnAddCard")?.addEventListener("click", async () => {
  const title = prompt("Nama card baru?");
  if (!title) return;

  await apiRequest("/cards", "POST", { title, color: "#3B82F6" });
  loadCards();
});

// ================================
// EDIT TITLE
// ================================
window.openEdit = (id) => {
  selectedCardId = id;
  const card = cards.find(c => c.id === id);
  const newTitle = prompt("Edit judul card:", card.title);
  if (!newTitle) return;

  apiRequest(`/cards/${id}`, "PUT", { title: newTitle, color: card.color })
    .then(loadCards)
    .catch(err => alert(err.message));
};

// ================================
// CHANGE COLOR
// ================================
window.openColor = async (id) => {
  selectedCardId = id;
  const card = cards.find(c => c.id === id);

  const color = prompt(
    "Masukkan HEX color (contoh #3B82F6)\nAtau pilih preset:\n" + presetColors.join(", "),
    card.color || "#3B82F6"
  );
  if (!color) return;

  await apiRequest(`/cards/${id}`, "PUT", { title: card.title, color });
  loadCards();
};

// ================================
// DELETE CARD
// ================================
window.deleteCard = async (id) => {
  if (!confirm("Yakin hapus card ini?")) return;
  await apiRequest(`/cards/${id}`, "DELETE");
  loadCards();
};

// init user info
const user = getUser();
document.getElementById("userEmail") && (document.getElementById("userEmail").innerText = user.email || "User");

loadCards();
