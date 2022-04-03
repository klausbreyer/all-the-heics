list = [];
// Setup the greet function
window.reload = function () {
  try {
    window.go.main.App.List()
      .then((result) => {
        // Update result with data back from App.List()
        list = result;
        refresh();
      })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);
  }
};

window.onload = function () {
  try {
    window.go.main.App.List()
      .then((result) => {
        // Update result with data back from App.List()
        list = result;
        refresh();
      })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);
  }
};

// Setup the greet function
window.greet = function () {
  for (let i = 0; i < list.length; i++) {
    list[i].started = true;
    refresh();
    try {
      window.go.main.App.WorkFile(list[i].Path)
        .then((result) => {
          list[i].worked = result;
          refresh();
        })
        .catch((err) => {
          console.error(err);
        });
    } catch (err) {
      console.error(err);
    }
  }
};

function refresh() {
  document.getElementById("result").innerHTML = list
    .map(
      (entry) =>
        `${entry.started && entry.worked == undefined ? "Started.." : ""}
        ${entry.worked ? "Finished:" : ""}
        ${entry.Path}
        ${entry.worked && entry.Correct ? "(Converted, Renamed, Deleted):" : ""}
        ${entry.worked && !entry.Correct ? "(Renamed):" : ""}
        <br />`
    )
    .join("");
}
