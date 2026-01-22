// frontend/js/sub_cards.js
import { apiRequest, getQueryParam, getUser } from "./api.js";

let subcards = [];
let selectedSubId = null;
const cardId = getQueryParam("card_id");

const presetColors = [
  "#3B82F6", "#8B5CF6", "#10B981", "#F97316",
  "#EF4444", "#14B8A6", "#0EA5E9", "#64748B"
];

async function loadSubCards() {
  if (!cardId) return alert("card_id tidak ditemukan!");
  subcards = await apiRequest(`/cards/${cardId}/subcards`);
  renderSubCards();
}

function renderSubCards() {
  const grid = document.getElementById("subGrid");
  grid.innerHTML = "";

  subcards.forEach(sc => {
    const el = document.createElement("div");
    el.className = "relative rounded-2xl p-5 text-white shadow-md hover:shadow-lg transition cursor-pointer";
    el.style.background = sc.color || "#8B5CF6";

    el.innerHTML = `
      <div class="flex items-start justify-between">
        <div>
          <h3 class="font-bold text-lg">${sc.title}</h3>
          <p class="text-xs opacity-90 mt-2">${sc.total_tasks || 0} Tasks • ${sc.completed || 0} Completed • ${sc.pending || 0} Pending</p>
        </div>

        <button class="menuBtn w-8 h-8 flex items-center justify-center rounded-full hover:bg-white/20 transition" data-id="${sc.id}">
          ⋮
        </button>
      </div>

      <div class="dropdown hidden absolute right-4 top-14 bg-white text-slate-800 rounded-xl shadow-lg w-44 overflow-hidden z-20" id="dropdown-${sc.id}">
        <button class="w-full text-left px-4 py-2 hover:bg-slate-100 text-sm" onclick="window.openEdit('${sc.id}')">Edit Judul</button>
        <button class="w-full text-left px-4 py-2 hover:bg-slate-100 text-sm" onclick="window.openColor('${sc.id}')">Ubah Warna</button>
        <button class="w-full text-left px-4 py-2 hover:bg-red-50 text-red-600 text-sm" onclick="window.deleteSubCard('${sc.id}')">Hapus</button>
      </div>
    `;

    el.addEventListener("click", () => {
      window.location.href = `tasks.html?subcard_id=${sc.id}`;
    });

    grid.appendChild(el);
  });

  // add box
  const addBox = document.createElement("div");
  addBox.className = "rounded-2xl border-2 border-dashed border-slate-300 bg-white flex flex-col items-center justify-center py-16 cursor-pointer hover:border-blue-400 transition";
  addBox.innerHTML = `
    <div class="w-12 h-12 rounded-full bg-slate-100 flex items-center justify-center text-xl font-bold text-slate-600">+</div>
    <p class="mt-3 font-semibold text-slate-700">Add Sub-Card</p>
  `;
  addBox.onclick = () => document.getElementById("btnAddSub").click();
  grid.appendChild(addBox);

  // menu
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

// Add Sub Card
document.getElementById("btnAddSub")?.addEventListener("click", async () => {
  const title = prompt("Nama sub-card baru?");
  if (!title) return;

  await apiRequest(`/cards/${cardId}/subcards`, "POST", { title, color: "#8B5CF6" });
  loadSubCards();
});

// Edit
window.openEdit = async (id) => {
  selectedSubId = id;
  const sc = subcards.find(s => s.id === id);
  const newTitle = prompt("Edit judul sub-card:", sc.title);
  if (!newTitle) return;

  await apiRequest(`/subcards/${id}`, "PUT", { title: newTitle, color: sc.color });
  loadSubCards();
};

// Color
window.openColor = async (id) => {
  selectedSubId = id;
  const sc = subcards.find(s => s.id === id);

  const color = prompt(
    "Masukkan HEX color (contoh #8B5CF6)\nPreset:\n" + presetColors.join(", "),
    sc.color || "#8B5CF6"
  );
  if (!color) return;

  await apiRequest(`/subcards/${id}`, "PUT", { title: sc.title, color });
  loadSubCards();
};

// Delete
window.deleteSubCard = async (id) => {
  if (!confirm("Yakin hapus sub-card ini?")) return;
  await apiRequest(`/subcards/${id}`, "DELETE");
  loadSubCards();
};

// user email
const user = getUser();
document.getElementById("userEmail") && (document.getElementById("userEmail").innerText = user.email || "User");

loadSubCards();
