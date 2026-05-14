// Sidebar toggle
function toggleSidebar() {
  var s = document.getElementById('sidebar');
  s.classList.toggle('hidden');
  var hidden = s.classList.contains('hidden');
  var edge = document.getElementById('edgeTrigger');
  var btn = document.getElementById('expandBtn');
  if (edge) edge.style.display = hidden ? 'block' : 'none';
  if (btn) btn.style.left = hidden ? '0' : '-40px';
}

function showSidebar() {
  var s = document.getElementById('sidebar');
  s.classList.remove('hidden');
  var edge = document.getElementById('edgeTrigger');
  var btn = document.getElementById('expandBtn');
  if (edge) edge.style.display = 'none';
  if (btn) btn.style.left = '-40px';
}

// Auth helpers
function getToken() {
  return localStorage.getItem('token');
}

async function fetchAPI(path, options) {
  options = options || {};
  var token = getToken();
  options.headers = options.headers || {};
  options.headers['Content-Type'] = 'application/json';
  if (token) options.headers['Authorization'] = 'Bearer ' + token;
  var res = await fetch(path, options);
  if (res.status === 401) {
    localStorage.removeItem('token');
    window.location.href = '/login';
    return;
  }
  return res.json();
}

// Toast
function showToast(msg, type) {
  var container = document.getElementById('toast-container');
  if (!container) return;
  var toast = document.createElement('div');
  toast.className = 'toast toast-' + (type || 'success');
  toast.textContent = msg;
  container.appendChild(toast);
  setTimeout(function() {
    toast.style.opacity = '0';
    setTimeout(function() { toast.remove(); }, 300);
  }, 3000);
}
