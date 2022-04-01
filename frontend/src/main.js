// Setup the greet function
window.greet = function () {
  document.getElementById("result").innerText =
    "Starting... (this can take a while)";
  try {
    window.go.main.App.Greet()
      .then((result) => {
        // Update result with data back from App.Greet()
        document.getElementById("result").innerText = result;
      })
      .catch((err) => {
        console.error(err);
      });
  } catch (err) {
    console.error(err);
  }
};

nameElement.onkeydown = function (e) {
  if (e.keyCode == 13) {
    window.greet();
  }
};
