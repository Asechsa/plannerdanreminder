// frontend/js/tasks.js
import { apiRequest, getQueryParam, getUser } from "./api.js";

let tasks = [];
let selectedTaskId = null;
const subcardId = getQueryParam("subcard_id");

// helpers
function formatDate(iso) {
  const d = new Date(iso);
  const month = d.toLocaleString("id-ID", { month: "short" });
  return `${d.getDate()} ${month}`;
}
function formatTime(iso) {
  const d = new Date(iso);
  return d.toLocaleTimeString("id-ID", { hour: "2-digit", minute: "2-digit" });
}
function badgeClass(level) {
  switch(level) {
    case "urgent": return "bg-red-100 text-red-600";
    case "overdue": return "bg-orange-100 text-orange-600";
    case "medium": return "bg-yellow-100 text-yellow-700";
    default: return "bg-slate-100 text-slate-600";
  }
}

function openModal(id) {
  const m = document.getElementById(id);
  m.classList.remove("hidden");
  m.classList.add("flex");
}
function closeModal(id) {
  const m = document.getElementById(id);
  m.classList.add("hidden");
  m.classList.remove("flex");
}
window.goBack = () => window.history.back();

// ======================================
// LOAD TASKS
// ======================================
async function loadTasks() {
  if (!subcardId) return alert("subcard_id tidak ditemukan!");
  tasks = await apiRequest(`/subcards/${subcardId}/tasks`);
  render();
}

function render() {
  const active = tasks.filter(t => t.status !== "done");
  const done = tasks.filter(t => t.status === "done");

  // stats
  document.getElementById("statTotal").textContent = tasks.length;
  document.getElementById("statCompleted").textContent = done.length;
  document.getElementById("statPending").textContent = active.length;

  document.getElementById("activeTitle").textContent = `Active Tasks (${active.length})`;
  document.getElementById("doneTitle").textContent = `Completed (${done.length})`;

  renderList("activeList", active, false);
  renderList("doneList", done, true);
}

function renderList(containerId, list, isDone) {
  const container = document.getElementById(containerId);
  container.innerHTML = "";

  list.forEach(task => {
    const row = document.createElement("div");
    row.className = "bg-white border rounded-2xl p-4 flex items-center justify-between shadow-sm";

    row.innerHTML = `
      <div class="flex items-center gap-4">
        <input type="checkbox" ${isDone ? "checked" : ""} class="w-5 h-5 accent-blue-600"/>
        <div>
          <p class="font-semibold ${isDone ? "line-through text-slate-400" : "text-slate-800"}">${task.title}</p>

          <div class="flex items-center gap-4 text-xs text-slate-500 mt-1">
            <span class="flex items-center gap-1">📅 ${formatDate(task.deadline_at)}</span>
            <span class="flex items-center gap-1">🕒 ${formatTime(task.deadline_at)}</span>
            <span class="px-2 py-0.5 rounded-full font-semibold ${badgeClass(task.urgency)}">${task.urgency}</span>
          </div>
        </div>
      </div>

      <div class="relative">
        <button class="menuBtn px-3 py-2 rounded-lg hover:bg-slate-100" data-id="${task.id}">⋮</button>

        <div class="dropdown hidden absolute right-0 top-10 bg-white border rounded-xl shadow-lg w-36 overflow-hidden z-20" id="drop-${task.id}">
          <button class="w-full text-left px-4 py-2 hover:bg-slate-100 text-sm" onclick="window.openEdit('${task.id}')">Edit</button>
          <button class="w-full text-left px-4 py-2 hover:bg-red-50 text-red-600 text-sm" onclick="window.deleteTask('${task.id}')">Delete</button>
        </div>
      </div>
    `;

    // toggle status
    row.querySelector("input").addEventListener("change", async () => {
      await apiRequest(`/tasks/${task.id}/status`, "PUT", {
        status: task.status === "done" ? "pending" : "done"
      });
      loadTasks();
    });

    container.appendChild(row);
  });

  // menu open
  document.querySelectorAll(".menuBtn").forEach(btn => {
    btn.onclick = (e) => {
      e.stopPropagation();
      toggleDropdown(btn.dataset.id);
    };
  });

  document.addEventListener("click", () => {
    document.querySelectorAll(".dropdown").forEach(d => d.classList.add("hidden"));
  });
}

function toggleDropdown(id) {
  document.querySelectorAll(".dropdown").forEach(d => d.classList.add("hidden"));
  document.getElementById("drop-" + id)?.classList.toggle("hidden");
}

// ======================================
// ADD TASK
// ======================================
document.getElementById("btnAddTask")?.addEventListener("click", () => openModal("modalAdd"));

document.getElementById("saveAddBtn")?.addEventListener("click", async () => {
  const title = document.getElementById("addTitle").value.trim();
  const deadline = document.getElementById("addDeadline").value;

  if (!title || !deadline) return alert("Judul & deadline wajib!");

  await apiRequest(`/subcards/${subcardId}/tasks`, "POST", {
    title,
    deadline_at: new Date(deadline).toISOString()
  });

  document.getElementById("addTitle").value = "";
  document.getElementById("addDeadline").value = "";
  closeModal("modalAdd");
  loadTasks();
});

// ======================================
// EDIT TASK (butuh endpoint PUT /api/tasks/:id)
// ======================================
window.openEdit = (id) => {
  selectedTaskId = id;
  const t = tasks.find(x => x.id === id);

  document.getElementById("editTitle").value = t.title;

  // ubah ISO -> datetime-local
  const dt = new Date(t.deadline_at);
  document.getElementById("editDeadline").value =
    dt.toISOString().slice(0,16);

  openModal("modalEdit");
};

document.getElementById("saveEditBtn")?.addEventListener("click", async () => {
  const title = document.getElementById("editTitle").value.trim();
  const deadline = document.getElementById("editDeadline").value;

  if (!title || !deadline) return alert("Judul & deadline wajib!");

  await apiRequest(`/tasks/${selectedTaskId}`, "PUT", {
    title,
    deadline_at: new Date(deadline).toISOString()
  });

  closeModal("modalEdit");
  loadTasks();
});

// ======================================
// DELETE TASK
// ======================================
window.deleteTask = async (id) => {
  if (!confirm("Yakin hapus task ini?")) return;
  await apiRequest(`/tasks/${id}`, "DELETE");
  loadTasks();
};

// user email
const user = getUser();
document.getElementById("userEmail") && (document.getElementById("userEmail").innerText = user.email || "User");

loadTasks();
