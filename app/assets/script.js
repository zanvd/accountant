document.querySelectorAll('.go-back').forEach(cancelButton => {
    cancelButton.addEventListener('click', () => window.history.back());
});

/* Deletion confirmation dialog. */
let deletionStorage = {
    submitUrl: ''
};
function setDeletionDialogContent(text, title) {
    document.getElementById('deletion-title').innerText = title;
    document.getElementById('deletion-text').innerHTML = text;
}
document.querySelectorAll('.confirm-deletion').forEach(element => element.addEventListener(
    'click',
    event => {
        event.preventDefault();
        event.stopPropagation();

        deletionStorage.submitUrl = element.getAttribute('href');

        let name = element.getAttribute('data-deletion-name');
        setDeletionDialogContent(
            `Are you sure you want to delete <strong>${name}</strong>?`,
            `Delete ${element.getAttribute('data-deletion-title')}`
        );

        document.getElementById('deletion-dialog-container').classList.toggle('hidden');
    }
));

document.getElementById('deletion-submit').addEventListener('click', event => {
    event.preventDefault();
    window.location.href = deletionStorage.submitUrl;
});

document.getElementById('deletion-cancel').addEventListener('click', event => {
    event.preventDefault();
    document.getElementById('deletion-dialog-container').classList.toggle('hidden');
    setDeletionDialogContent('', '');
    deletionStorage.submitUrl = '';
});
/* Deletion confirmation dialog end. */
