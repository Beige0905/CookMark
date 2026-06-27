let dark = $state(false);

function init() {
	dark = localStorage.getItem('theme') === 'dark';
}

function toggle() {
	dark = !dark;
	localStorage.setItem('theme', dark ? 'dark' : 'light');
	document.documentElement.classList.toggle('dark', dark);
}

export const theme = { get dark() { return dark; }, init, toggle };
