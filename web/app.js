const output = document.querySelector("#output");

function show(value) {
  output.textContent = JSON.stringify(value, null, 2);
}

async function request(url, options) {
  const response = await fetch(url, options);
  const data = await response.json().catch(() => response.text());
  if (!response.ok) throw new Error(JSON.stringify(data));
  show(data);
}

function noteBody() {
  return JSON.stringify({
    title: document.querySelector("#note-title").value,
    content: document.querySelector("#note-content").value,
  });
}

document.querySelector("#create-note").onclick = async () => {
  try {
    const response = await fetch("/notes", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: noteBody(),
    });
    const data = await response.json();
    if (!response.ok) throw new Error(JSON.stringify(data));

    document.querySelector("#note-id").value = data.id;
    show(data);
  } catch (error) {
    show({ error: error.message });
  }
};

document.querySelector("#list-notes").onclick = () => request("/notes")
  .catch(error => show({ error: error.message }));

document.querySelector("#get-note").onclick = () => request(
  `/notes/${document.querySelector("#note-id").value}`,
).catch(error => show({ error: error.message }));

document.querySelector("#update-note").onclick = () => request(
  `/notes/${document.querySelector("#note-id").value}`,
  {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: noteBody(),
  },
).catch(error => show({ error: error.message }));

document.querySelector("#delete-note").onclick = () => request(
  `/notes/${document.querySelector("#note-id").value}`,
  { method: "DELETE" },
).catch(error => show({ error: error.message }));

document.querySelector("#upload-file").onclick = () => {
  const form = new FormData();
  form.append("file", document.querySelector("#file").files[0]);
  request("/files", { method: "POST", body: form })
    .catch(error => show({ error: error.message }));
};

document.querySelector("#list-files").onclick = () => request("/files")
  .catch(error => show({ error: error.message }));

document.querySelector("#view-file").onclick = async () => {
  const key = document.querySelector("#file-key").value;
  try {
    const response = await fetch(`/files/${encodeURIComponent(key)}`);
    if (!response.ok) throw new Error(await response.text());
    show(await response.text());
  } catch (error) {
    show({ error: error.message });
  }
};
