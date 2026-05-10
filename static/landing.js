const button = document.getElementById("create");
const loading = document.getElementById("loading");
const toast = document.getElementById("toast");
const toastBootstrap = bootstrap.Toast.getOrCreateInstance(toast);

function handleError() {
    toastBootstrap.show();
    button.style.display = "";
    loading.style.display = "none";
}

button.addEventListener("click", async (e) => {
    button.style.display = "none";
    loading.style.display = "";
    try {
        const response = await fetch("/api/document", {
            method: "POST"
        });
        if (!response.ok) {
            handleError();
            return;
        }

        const body = await response.json();
        window.location.href = `/document/${body.id}/edit`;
    } catch(error) {
        console.error(error);
        handleError();
    }
});
