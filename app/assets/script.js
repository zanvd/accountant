Array.from(document.getElementsByClassName('btn-cancel')).forEach(cancelButton => {
    cancelButton.addEventListener('click', () => window.location = '/');
});
