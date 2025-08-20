// Получить текущее состояние поля (массив из 9 булевых значений)
async function fetchState() {
  const res = await fetch('api/state');
  return res.json();
}

// Переключить клетку по индексу (1..9): меняется выбранная и её соседи сверху/снизу/слева/справа
async function toggle(index) {
  const params = new URLSearchParams();
  params.set('index', String(index));
  const res = await fetch('api/toggle', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: params.toString()
  });
  return res.json();
}

// Сбросить поле в случайное состояние
async function resetBoard() {
  const res = await fetch('api/reset', { method: 'POST' });
  return res.json();
}

// Текущее состояние и ссылки на DOM-элементы сфер
let state = [];
const spheres = [];

// Применить состояние к DOM: навесить классы white/black
function applyState() {
  for (let i = 0; i < 9; i++) {
    const el = spheres[i];
    if (!el) continue;
    const isWhite = !!state[i];
    el.classList.toggle('white', isWhite);
    el.classList.toggle('black', !isWhite);
    el.setAttribute('aria-pressed', isWhite ? 'true' : 'false');
  }
}

document.addEventListener('DOMContentLoaded', async () => {
  // Инициализация ссылок и обработчиков клика по сферам
  document.querySelectorAll('.sphere').forEach((btn) => {
    const idx = Number(btn.getAttribute('data-index'));
    spheres[idx - 1] = btn;
    // До загрузки состояния помечаем как белые для первого рендера
    btn.classList.add('white');
    btn.addEventListener('click', async () => {
      state = await toggle(idx);
      applyState();
    });
  });

  // Кнопка «Случайное поле»
  document.getElementById('reset').addEventListener('click', async () => {
    state = await resetBoard();
    applyState();
  });

  // Клавиатура: 1..9 (и NumPad 1..9) переключают соответствующие клетки
  document.addEventListener('keydown', async (e) => {
    const key = e.key;
    if (/^[1-9]$/.test(key)) {
      e.preventDefault();
      const idx = Number(key);
      state = await toggle(idx);
      applyState();
    }
    // Поддержка цифрового блока (numpad)
    if (e.code && /^Numpad[1-9]$/.test(e.code)) {
      e.preventDefault();
      const idx = Number(e.code.replace('Numpad', ''));
      state = await toggle(idx);
      applyState();
    }
  });

  // Первичная загрузка состояния
  state = await fetchState();
  applyState();
});


