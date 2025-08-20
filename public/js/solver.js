// Управление состоянием локальной конфигурации на странице решателя
const sSpheres = [];
let sState = Array(9).fill(true); // по умолчанию все белые

function renderSolverState() {
  for (let i = 0; i < 9; i++) {
    const el = sSpheres[i];
    if (!el) continue;
    const isWhite = !!sState[i];
    el.classList.toggle('white', isWhite);
    el.classList.toggle('black', !isWhite);
    el.setAttribute('aria-pressed', isWhite ? 'true' : 'false');
  }
}

function toBitString() {
  return sState.map(v => (v ? '1' : '0')).join('');
}

async function solve(target) {
  const params = new URLSearchParams();
  params.set('target', target);
  params.set('state', toBitString());
  const res = await fetch('api/solve', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: params.toString()
  });
  return res.json();
}

async function fetchCurrentState() {
  const res = await fetch('api/state');
  return res.json();
}

document.addEventListener('DOMContentLoaded', () => {
  document.querySelectorAll('.sphere').forEach(btn => {
    const idx = Number(btn.getAttribute('data-index')) - 1;
    sSpheres[idx] = btn;
    btn.addEventListener('click', () => {
      sState[idx] = !sState[idx];
      renderSolverState();
    });
  });

  document.getElementById('solve').addEventListener('click', async () => {
    const target = (document.querySelector('input[name="target"]:checked')?.value) || 'white';
    const result = await solve(target);
    const el = document.getElementById('solution');
    const seq = (result.moves || []).join(', ');
    // Показываем только последовательность ходов
    el.innerHTML = `Ходы: <code>${seq}</code>`;
  });

  document.getElementById('fill-from-current').addEventListener('click', async () => {
    const current = await fetchCurrentState();
    // current — массив true/false
    sState = Array.from(current).map(Boolean);
    renderSolverState();
  });

  document.getElementById('clear').addEventListener('click', () => {
    sState = Array(9).fill(false);
    renderSolverState();
  });

  // Первичная отрисовка
  renderSolverState();
});


