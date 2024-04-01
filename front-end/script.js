document.addEventListener('DOMContentLoaded', function() {
    navigateToPage(1); // Load default page
});

function navigateToPage(pageNumber) {
    fetchPage(pageNumber)
        .then(response => response.text())
        .then(html => {
            document.getElementById('content').innerHTML = html;
        })
        .catch(error => console.error('Error loading page: ', error));
}

function fetchPage(pageNumber) {
    let url = '';
    switch (pageNumber) {
        case 1:
            url = 'mathExpressions.html';
            break;
        case 2:
            url = 'agentSettings.html';
            break;
        case 3:
            url = 'agentMonitoring.html';
            break;
        case 4:
            url = 'login.html';
            break;
        case 5:
            url = 'register.html';
            break;
        default:
            url = 'mathExpressions.html';
            break;
    }
    return fetch(url);
}
