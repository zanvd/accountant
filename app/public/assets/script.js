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

/* Dashboard cards. */
document.querySelectorAll('.dashboard-card').forEach(element => element.addEventListener(
    'click',
    _ => {
        let name = element.getAttribute('data-transaction-name');
        let category = element.getAttribute('data-transaction-category');
        let type = element.getAttribute('data-transaction-type');
        window.location.href = `/app/transaction/add?name=${name}&category=${category}&type=${type}`;
    }
));
/* Dashboard cards end. */

/* Transactions filter. */
document.getElementById('transactions-filter-submit').addEventListener('click', function () {
    let params = new URLSearchParams(window.location.search);

    let catId = document.getElementById('transactions-filter-category').value;
    if (catId) {
        params.set('category', catId);
    } else {
        params.delete('category')
    }

    window.location.search = params.toString();
});

Array.from(document.getElementsByClassName('transactions-month-button'), link => link.addEventListener('click', function () {
    let linkParams = new URL(link.href).searchParams;
    let url = new URL(window.location);

    url.searchParams.set('from', linkParams.get('from'));
    url.searchParams.set('to', linkParams.get('to'));

    link.href = url.href;
}));
/* Transactions filter end. */
