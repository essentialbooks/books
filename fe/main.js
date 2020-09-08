import Toc from "./Toc.svelte";
import SearchInput from "./SearchInput.svelte";
import PageToc from "./PageToc.svelte";
import BookToc from "./BookToc.svelte";
import { item } from "./item.js";
import { viewGet, viewSet, viewClear } from "./store.js";
import { logflare } from "./util.js";

// pageId looks like "5ab3b56329c44058b5b24d3f364183ce"
// find full url of the page matching this pageId
function findURLWithPageId(pageId) {
  var n = gTocItems.length;
  for (var i = 0; i < n; i++) {
    var tocItem = gTocItems[i];
    var uri = item.url(tocItem);
    // uri looks like "go-get-5ab3b56329c44058b5b24d3f364183ce"
    if (uri.endsWith(pageId)) {
      return uri;
    }
  }
  return "";
}

function do404() {
  var loc = window.location.pathname;
  var locParts = loc.split("/");
  var lastIdx = locParts.length - 1;
  var uri = locParts[lastIdx];
  // redirect ${garbage}-${id} => ${correct url}-${id}
  var parts = uri.split("-");
  var pageId = parts[parts.length - 1];
  var fullURL = findURLWithPageId(pageId);
  if (fullURL != "") {
    locParts[lastIdx] = fullURL;
    var loc = locParts.join("/");
    window.location.pathname = loc;
  }
}
window.do404 = do404;

function httpsMaybeRedirect() {
  if (window.location.protocol !== "http:") {
    return;
  }
  if (window.location.hostname !== "www.programming-books.io") {
    return;
  }
  var uri = window.location.toString();
  uri = uri.replace("http://", "https://");
  window.location = uri;
}

window.httpsMaybeRedirect = httpsMaybeRedirect;

function showContact() {
  var el = document.getElementById("contact-form");
  el.style.display = "block";
  el = document.getElementById("contact-page-url");
  var uri = window.location.href;
  //uri = uri.replace("#", "");
  el.value = uri;
  el = document.getElementById("msg-for-chris");
  el.focus();
}

function hideContact() {
  var el = document.getElementById("contact-form");
  el.style.display = "none";
}

window.showContact = showContact;
window.hideContact = hideContact;

// rv = rememberView but short because it's part of url
function rv(view) {
  //console.log("rv:", view);
  viewSet(view);
}
window.rv = rv;

const app = {
  toc: Toc,
  searchInput: SearchInput,
  pageToc: PageToc,
  bookToc: BookToc,
};

function updateLinkHome() {
  var view = viewGet();
  if (!view) {
    return;
  }
  var uri = "/";
  if (view === "list") {
    // do nothing
  } else if (view == "grid") {
    uri = "/index-grid";
  } else {
    console.log("unknown view:", view);
    viewClear();
  }
  var el = document.getElementById("link-home");
  if (el && el.href) {
    //console.log("update home url to:", uri);
    el.href = uri;
  }
}

function doIndexPage() {
  var view = viewGet();
  var loc = window.location.pathname;
  //console.log("doIndexPage(): view:", view, "loc:", loc);
  if (!view) {
    return;
  }
  if (view === "list") {
    if (loc === "/index-grid") {
      window.location = "/";
    }
  } else if (view === "grid") {
    if (loc === "/") {
      window.location = "/index-grid";
    }
  } else {
    console.log("Unknown view:", view);
  }
}

function logCurrentURL() {
  const path = window.location.pathname;
  const host = window.location.hostname;
  const book = host.split(".")[0];
  const d = {
    "log_entry": `page ${book} ${host}${path}`,
    meatadata: {
      url: path,
      host: host,
      book: book,
    },
  }
  logflare(d);
}

logCurrentURL();

// we don't want to run javascript on about etc. pages
var loc = window.location.pathname;
var isIndexPage = loc === "/" || loc === "/index-grid";

if (isIndexPage) {
  doIndexPage();
  updateLinkHome();
}

export default app;
