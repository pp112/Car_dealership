const form = document.getElementById('searchForm');
const fieldSelect = document.getElementById('field');
const numericSearch = document.getElementById('numericSearch');
const textLabel = document.getElementById('textLabel');
const qInput = document.getElementById('q');

// ------------------------ SUBMIT SEARCH ------------------------

form.addEventListener('submit', async (e) => {
  e.preventDefault();

  const field = fieldSelect.value;
  let q = '';

  // Числовые поля: year, price
  if (field === 'year' || field === 'price') {
    const op = document.getElementById('operator').value;
    const val = document.getElementById('numericValue').value.trim();

    if (!val) {
      alert('Введите число');
      return;
    }

    q = op + val; // Например ">2015" или "<=1200"
  } else {
    // Строковые поля: brand, model
    q = qInput.value.trim();
    if (!q) {
      alert('Введите поисковый запрос');
      return;
    }
  }

  const url = `/api/search?field=${encodeURIComponent(field)}&q=${encodeURIComponent(q)}`;

  sendRequest(url);
});

// ------------------------ SHOW ALL ------------------------

document.getElementById('showAllBtn').addEventListener('click', () => {
  sendRequest('/api/all');
});

// ------------------------ SEND REQUEST + RENDER ------------------------

async function sendRequest(url) {
  try {
    const resp = await fetch(url);

    if (!resp.ok) {
      const err = await resp.json();
      alert('Ошибка: ' + (err.error || resp.statusText));
      return;
    }

    const data = await resp.json();
    const table = document.getElementById('resTable');
    const body = document.getElementById('resBody');
    const nores = document.getElementById('nores');

    body.innerHTML = '';

    if (Array.isArray(data) && data.length > 0) {
      data.forEach(row => {
        const tr = document.createElement('tr');

        tr.innerHTML = `
          <td>${row.id}</td>
          <td>${row.brand}</td>
          <td>${row.model}</td>
          <td>${row.year}</td>
          <td>${row.price}</td>
        `;

        body.appendChild(tr);
      });

      table.style.display = '';
      nores.style.display = 'none';
    } else {
      table.style.display = 'none';
      nores.style.display = 'block';
    }

  } catch (err) {
    alert('Fetch error: ' + err);
  }
}

// ------------------------ SWITCH FIELD ------------------------

fieldSelect.addEventListener('change', () => {
  const field = fieldSelect.value;

  if (field === 'year' || field === 'price') {
    numericSearch.style.display = 'block';
    textLabel.style.display = 'none';
    qInput.style.display = 'none';
  } else {
    numericSearch.style.display = 'none';
    textLabel.style.display = 'block';
    qInput.style.display = 'block';
  }
});
