// import queue from "./response_queue_1765768521331_mine.json" with { type: "json" };
import queue from "./response_queue_1765767430744_erica.json" with { type: "json" };

// playlistPanelVideoRenderer
// playlistPanelVideoWrapperRenderer

let s = "Now: ";
let n = 0;
let foundSelected = false;
for (const v of queue.items) {
  if (v.playlistPanelVideoWrapperRenderer) {
    (v.playlistPanelVideoRenderer as any) =
      v.playlistPanelVideoWrapperRenderer.primaryRenderer.playlistPanelVideoRenderer;
  }
  if (!v.playlistPanelVideoRenderer) {
    throw new Error("fkd up");
  }
  if (v.playlistPanelVideoRenderer.selected) {
    foundSelected = true;
  }
  if (!v.playlistPanelVideoRenderer.selected && !foundSelected) {
    continue;
  }
  if (n > 5) {
    break;
  }
  n++;
  const title = v.playlistPanelVideoRenderer.title.runs[0].text;
  const artist = v.playlistPanelVideoRenderer.shortBylineText.runs[0].text;
  // const artist = v.playlistPanelVideoRenderer?.shortBylineText.runs[0].text
  let sl = "";
  if (n === 1) {
    sl = title + " - " + artist + ", ";
  } else {
    sl = "#" + (n - 1) + ": " + title + " - " + artist + ", ";
  }
  s += sl;
}
console.log(s);
