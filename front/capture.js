const KEY = 'gtd-capture';
const $input = document.getElementById('input');
const $list = document.getElementById('list');
const $empty = document.getElementById('empty');

// rawStuffs without id
let rawStuffs = [];

let stuffs = [];
try { stuffs = JSON.parse(localStorage.getItem(KEY)) || []; } catch (e) { }

const save = () => localStorage.setItem(KEY, JSON.stringify(stuffs));

function render() {
    $list.innerHTML = '';
    stuffs.forEach((text, i) => {
        const li = document.createElement('li');
        const span = document.createElement('span');
        span.textContent = text;

        const btn = document.createElement('button');
        btn.type = 'button';
        btn.textContent = '×';
        btn.title = 'delete';
        btn.onclick = () => { stuffs.splice(i, 1); save(); render(); };

        li.append(span, btn);
        $list.append(li);
    });
    $empty.style.display = stuffs.length ? 'none' : 'block';
}

// enter to capture
document.getElementById('form').addEventListener('submit', e => {
    e.preventDefault();

    // try to send request to server, then get id from server
    const text = $input.value.trim();
    if (!text) return;
    stuffs.unshift(text);   // 最新的在最上面
    save();
    $input.value = '';
    render();
    $input.focus();
});

render();
$input.focus();
