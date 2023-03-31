const addBasketButtons = document.querySelectorAll('button[type="addbasket"]');
const basketCount = document.querySelector('.basket-count');

let count = 0;

addBasketButtons.forEach(button => {
  button.addEventListener('click', () => {
    count++;
    basketCount.textContent = count;
  });
});
