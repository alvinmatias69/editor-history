let isStop = false;
let stopIdx = 0;
let shouldUpdateSlider = false;
let currentDocumentID;

const susOperationMaxLen = 25;

const finishBtn = document.getElementById("finish-btn");
const progressBtn = document.getElementById("on-progress-btn");
const editor = document.getElementById("editor");
const slider = document.getElementById("slider");
const play = document.getElementById("play");
const pause = document.getElementById("pause");
const back = document.getElementById("back");
const forward = document.getElementById("forward");
const readUrlBtn = document.getElementById("read-url-btn");
const readUrl = document.getElementById("read-url");
const errorModal = new bootstrap.Modal("#error-modal");
const retryBtn = document.getElementById("retry-btn");
const suspiciousEmpty = document.getElementById("suspicious-empty");
const suspiciousParent = document.getElementById("suspicious-parent");
const suspiciousOffcanvas = new bootstrap.Offcanvas("#suspicious-list");

const OPERATION_TYPE = {
    INSERT: "insert",
    DELETE: "delete",
};

const defaultConfig = {
    clearEditor: false,
    delayMs: 0,
    shouldUpdateSlider: false,
    shouldUpdateEditorPerState: false,
    stopIdx: 0,
    startIdx: 0,
}

async function init(documentID) {
    currentDocumentID = documentID;
    editor.value = "";
    let response;

    try {
        response = await fetch(`/api/document/${documentID}/read`, {
            method: "GET"
        });

        if (!response.ok) {
            errorModal.show();
            return;
        }
    } catch(err) {
        errorModal.show();
        return;
    }

    const body = await response.json();

    if (body.is_finished) {
        finishBtn.style.display = "";
    } else {
        progressBtn.style.display = "";
    }

    operations = body.operations;
    slider.max = operations.length - 1;
    const config = structuredClone(defaultConfig);
    config.stopIdx = operations.length - 1;
    config.clearEditor = true;
    await constructEditor(config);
    slider.value = operations.length - 1;

    setShareUrl(documentID);
    hydrateSuspicious();
}

function sanitize(word) {
    const sanitized = word.replaceAll("\n", " ");

    return sanitized.length > susOperationMaxLen
        ? sanitized.substring(0, susOperationMaxLen) + "..."
        : sanitized;
}

function hydrateSuspicious() {
    const list = [];
    for (const idx in operations) {
        const operation = operations[idx];
        if (operation.value.length < 2) {
            continue;
        }

        const subheading = document.createElement("div");
        subheading.setAttribute("class", "fw-bold text-truncate");
        subheading.innerText = sanitize(operation.value);

        const operationDate = new Date(operation.operation_timestamp);

        const timestamp = document.createElement("small");
        timestamp.setAttribute("class", "text-secondary fst-italic");
        timestamp.innerText = `${operationDate.toLocaleDateString()} ${operationDate.toLocaleTimeString()}`;

        const container = document.createElement("div");
        container.setAttribute("class", "ms-2 me-auto");
        container.append(subheading, timestamp);

        const count = document.createElement("span");
        count.setAttribute("class", "badge text-bg-warning rounded-pill");
        count.innerText = operation.value.length;

        const button = document.createElement("button");
        button.setAttribute("type", "button");
        button.setAttribute("class", "list-group-item list-group-item-action d-flex justify-content-between align-items-start");
        button.append(container, count);
        button.addEventListener("click", event => {
            slider.value = idx;
            handleSliderUpdate();
            suspiciousOffcanvas.hide();
        });
        list.push(button);
    }

    if (list.length === 0) {
        return;
    }

    suspiciousEmpty.style.display = "none";
    suspiciousParent.append(...list);
}

async function constructEditor(config) {
    let value = config.clearEditor ? "" : editor.value;
    let idx = config.startIdx
    while (idx <= config.stopIdx) {
        const event = operations[idx];
        if (event.operation == OPERATION_TYPE.INSERT) {
            value = value.slice(0, event.position_start) + event.value + value.slice(event.position_end, value.length);
        } else if (event.operation == OPERATION_TYPE.DELETE) {
            const start = event.position_start === event.position_end
                  ? event.position_start - 1
                  : event.position_start;
            value = value.slice(0, start) + value.slice(event.position_end, value.length);
        }

        if (config.shouldUpdateEditorPerstate) {
            editor.value = value;
        }

        if (config.shouldUpdateSlider) {
            slider.value = idx;
        }

        await sleep(config.delayMs);
        idx++;

        if (isStop) {
            break;
        }
    }

    editor.value = value;
    play.style.display = "";
    pause.style.display = "none";
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

function handleSliderUpdate() {
    isStop = false;
    const config = structuredClone(defaultConfig);
    config.stopIdx = parseInt(slider.value);
    config.shouldUpdateSlider = false;
    config.shouldUpdateEditorPerstate = false;
    config.clearEditor = true;
    constructEditor(config);
}

function setShareUrl(documentID) {
    const baseURL = window.location.origin;
    readUrl.value = `${baseURL}/document/${documentID}/view`;
}

play.addEventListener("click", () => {
    const sliderValue = parseInt(slider.value);
    play.style.display = "none";
    pause.style.display = "";
    isStop = false;

    const config = structuredClone(defaultConfig);
    config.clearEditor = sliderValue === 0 || sliderValue === operations.length - 1;
    config.delayMs = 100;
    config.shouldUpdateSlider = true;
    config.shouldUpdateEditorPerstate = true;
    config.stopIdx = operations.length - 1;
    if (config.clearEditor) {
        config.startIdx = sliderValue;
    } else {
        config.startIdx = sliderValue === operations.length - 1
            ? 0
            : sliderValue + 1;
    }
    if (sliderValue == operations.length - 1) {
        config.startIdx = 0;
    }

    constructEditor(config);
});

pause.addEventListener("click", () => {
    isStop = true;
    play.style.display = "";
    pause.style.display = "none";
});

slider.addEventListener("change", e => handleSliderUpdate());

back.addEventListener("click", () => {
    const sliderValue = parseInt(slider.value);
    if (sliderValue === 0) {
        return;
    }
    slider.value = sliderValue - 1;
    handleSliderUpdate();
});

forward.addEventListener("click", () => {
    const sliderValue = parseInt(slider.value);
    if (sliderValue === (operations.length - 1)) {
        return;
    }
    slider.value = sliderValue + 1;
    handleSliderUpdate();
});

readUrlBtn.addEventListener("click", e => {
    navigator.clipboard.writeText(readUrl.value);
});

retryBtn.addEventListener("click", e => {
    init(currentDocumentID)
});
