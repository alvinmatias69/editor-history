const SESSION_ID_KEY = "session_id";

const OPERATION_TYPE = {
    INSERT: "insert",
    DELETE: "delete",
};

const EVENT_MAP = {
    "insertText": OPERATION_TYPE.INSERT,
    "insertLineBreak": OPERATION_TYPE.INSERT,
    "insertParagraph": OPERATION_TYPE.INSERT,
    "insertFromPaste": OPERATION_TYPE.INSERT,
    "deleteWordBackward": OPERATION_TYPE.DELETE,
    "deleteWordForward": OPERATION_TYPE.DELETE,
    "deleteContentBackward": OPERATION_TYPE.DELETE,
    "deleteContentForward": OPERATION_TYPE.DELETE,
}

const editor = document.getElementById("editor");
const fullSyncButton = document.getElementById("full-sync");
const readUrl = document.getElementById("read-url");
const writeUrl = document.getElementById("write-url");
const readUrlBtn = document.getElementById("read-url-btn");
const writeUrlBtn = document.getElementById("write-url-btn");
const lastSyncTimestamp = document.getElementById("last-sync-timestamp");
const finish = document.getElementById("finish");
const finishBtn = document.getElementById("finish-btn");
const retryBtn = document.getElementById("retry-btn");
const overrideBtn = document.getElementById("override-btn");

const accessModal = new bootstrap.Modal("#access-modal");
const notFoundModal = new bootstrap.Modal("#not-found-modal");
const errorModal = new bootstrap.Modal("#error-modal");
const finishModal = new bootstrap.Modal("#finish-modal");
const finishErrorModal = new bootstrap.Modal("#finish-error-modal");

let retryCb = () => {};
let accessCb = () => {};

let positionStart = 0;
let positionEnd = 0;
let documentID;
let readerID;
let operations = [];
let timeoutID;
let lastSendID = null;

function updatePosition() {
    positionStart = editor.selectionStart;
    positionEnd = editor.selectionEnd;
}

function getSessionID() {
    let sessionID = localStorage.getItem(SESSION_ID_KEY);
    if (!!sessionID) {
        return sessionID;
    }

    sessionID = crypto.randomUUID();
    localStorage.setItem(SESSION_ID_KEY, sessionID);
    return sessionID;
}

function setShareUrl() {
    const baseURL = window.location.origin;
    readUrl.value = `${baseURL}/document/${readerID}/view`;
    writeUrl.value = `${baseURL}/document/${documentID}/edit`;
}

async function init(initialDocumentID) {
    documentID = initialDocumentID;

    const params = new URLSearchParams();
    params.append("editor-id", getSessionID());

    let response;
    try {
        response = await fetch(`/api/document/${documentID}/write?${params}`, {
            method: "GET"
        });
    } catch (err) {
        handleError(response, () => init(initialDocumentID));
        return;
    }

    if (!response.ok) {
        handleError(response, () => init(initialDocumentID));
        return;
    }

    const body = await response.json();
    readerID = body.reader_id;
    operations = body.operations;
    constructEditor();
    if (operations.length > 0) {
        lastSendID = operations[operations.length-1].id;
    }
    setShareUrl();
    updateLastSync();
    if (body.is_finished) {
        setFinishState();
    }
}

async function handleFinish() {
    const payload = {
        editor_id: getSessionID(),
    };

    let response;
    try {
        response = await fetch(`/api/document/${documentID}/finish`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });
    } catch (err) {
        handleError(response, handleFinish);
        return;
    }

    if (!response.ok) {
        handleError(response, handleFinish);
        return;
    }

    setFinishState();
}

function setFinishState() {
    editor.readOnly = true;
    editor.disabled = true;
    finishBtn.disabled = true;
    finishBtn.innerText = "Finished ✓";
    fullSyncButton.disabled = true;
    finishModal.hide();
}

function fullSync() {
    sendOperations(0);
}

function partialSync() {
    let startIdx = 0;
    if (!!lastSendID) {
        startIdx = operations.findIndex(item => item.id == lastSendID) + 1;
    }
    sendOperations(startIdx);
}

async function sendOperations(startIdx) {
    const operationRequests = operations.slice(startIdx);

    const payload = {
        editor_id: getSessionID(),
        operation_requests: operationRequests,
    }

    let response;
    try {
        response = await fetch(`/api/document/${documentID}/operations`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });
    } catch (err) {
        handleError(response, () => sendOperations(startIdx));
        return;
    }


    if (!response.ok) {
        handleError(response, () => sendOperations(startIdx));
        return;
    }

    if (operationRequests.length > 0) {
        lastSendID = operationRequests[operationRequests.length-1].id;
    }
    updateLastSync();
}

function debounce() {
    clearTimeout(timeoutID);
    timeoutID = setTimeout(() => {
        partialSync();
    }, 1000);
}

function constructEditor() {
    let value = "";
    for (event of operations) {
        if (event.operation == OPERATION_TYPE.INSERT) {
            value = value.slice(0, event.position_start) + event.value + value.slice(event.position_end, value.length);
        } else if (event.operation == OPERATION_TYPE.DELETE) {
            const start = event.position_start === event.position_end
                  ? event.position_start - 1
                  : event.position_start;
            value = value.slice(0, start) + value.slice(event.position_end, value.length);
        }
    }
    editor.value = value;
    autogrow();
}

function autogrow() {
    if (editor.scrollHeight > editor.clientHeight) {
        editor.style.height = `${editor.scrollHeight}px`;
    }
}

function updateLastSync() {
    const now = (new Date()).toString();
    lastSyncTimestamp.innerText = now;
}

function handleRetryError(callback) {
    retryBtn.removeEventListener("click", retryCb);
    retryCb = () => {
        callback();
    };
    retryBtn.addEventListener("click", retryCb);
    errorModal.show();
}

async function handleError(response, callback) {
    if (!response) {
        handleRetryError(callback);
        return;
    }

    if (response.status == 404) {
        notFoundModal.show();
        return;
    }

    if (response.status == 409) {
        const body = await response.json();
        if (body.code === 2) {
            finishErrorModal.show();
            init(documentID);
            return;
        }

        if (body.code === 3) {
            accessModal.show();
            overrideBtn.removeEventListener("click", accessCb);
            accessCb = () => {
                overrideAccess(callback);
            };
            overrideBtn.addEventListener("click", accessCb);
            return;
        }
    }

    handleRetryError(callback);
}

async function overrideAccess(callback) {
    const payload = {
        editor_id: getSessionID(),
    }

    let response;
    try {
        response = await fetch(`/api/document/${documentID}/override`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });
    } catch (err) {
        handleError(response, () => overrideAccess(callback));
        return;
    }


    if (!response.ok) {
        handleError(response, () => overrideAccess(callback));
        return;
    }

    callback();
}

editor.addEventListener("keyup", e => {
    updatePosition();
    autogrow();
});

editor.addEventListener("selectionchange", e => {
    updatePosition();
});

editor.addEventListener("input", e => {
    let value = e.data;
    if (e.inputType == "insertLineBreak") {
        value = "\n";
    }

    const data = {
        id: crypto.randomUUID(),
        editor_id: getSessionID(),
        value: value,
        position_start: positionStart,
        position_end: positionEnd,
        operation: EVENT_MAP[e.inputType],
        timestamp: (new Date()).toISOString(),
    };

    operations.push(data);
    debounce();
});

fullSyncButton.addEventListener("click", e => {
    fullSync();
});

readUrlBtn.addEventListener("click", e => {
    navigator.clipboard.writeText(readUrl.value);
});

writeUrlBtn.addEventListener("click", e => {
    navigator.clipboard.writeText(writeUrl.value);
});

finish.addEventListener("click", e => {
    handleFinish();
});
