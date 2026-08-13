const API = "";

export async function api(path, opts) {
  const res = await fetch(API + path, opts);
  const data = await res.json().catch(() => null);
  if (!res.ok) throw new Error((data && data.error) || "Request failed");
  return data;
}