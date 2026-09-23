const output = document.querySelector("#output");

function show(value) {
  output.textContent = JSON.stringify(value, null, 2);
}

async function request(url, options) {
  const response = await fetch(url, options);
  const data = await response.json().catch(() => response.text());
  if (!response.ok) throw new Error(JSON.stringify(data));
  show(data);
  return data;
}

function noteBody(titleId, contentId) {
  return JSON.stringify({
    title: document.querySelector(`#${titleId}`).value,
    content: document.querySelector(`#${contentId}`).value,
  });
}

document.querySelector("#create-note").onclick = async () => {
  try {
    const note = await request("/notes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: noteBody("create-title", "create-content"),
    });
    document.querySelector("#get-note-id").value = note.id;
    document.querySelector("#update-note-id").value = note.id;
    document.querySelector("#delete-note-id").value = note.id;
  } catch (error) {
    show({ error: error.message });
  }
};

document.querySelector("#list-notes").onclick = () => request("/notes")
  .catch(error => show({ error: error.message }));

document.querySelector("#get-note").onclick = () => request(
  `/notes/${document.querySelector("#get-note-id").value}`,
).catch(error => show({ error: error.message }));

document.querySelector("#update-note").onclick = () => request(
  `/notes/${document.querySelector("#update-note-id").value}`,
  {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: noteBody("update-title", "update-content"),
  },
).catch(error => show({ error: error.message }));

document.querySelector("#delete-note").onclick = () => request(
  `/notes/${document.querySelector("#delete-note-id").value}`,
  { method: "DELETE" },
).catch(error => show({ error: error.message }));

document.querySelector("#upload-file").onclick = () => {
  const selectedFile = document.querySelector("#file").files[0];
  if (!selectedFile) {
    show({ error: "choose a file first" });
    return;
  }

  const form = new FormData();
  form.append("file", selectedFile);
  request("/files", { method: "POST", body: form })
    .catch(error => show({ error: error.message }));
};

document.querySelector("#list-files").onclick = () => request("/files")
  .catch(error => show({ error: error.message }));

document.querySelector("#view-file").onclick = async () => {
  const key = document.querySelector("#view-file-key").value;
  try {
    const response = await fetch(`/files/${encodeURIComponent(key)}`);
    if (!response.ok) throw new Error(await response.text());
    show(await response.text());
  } catch (error) {
    show({ error: error.message });
  }
};
