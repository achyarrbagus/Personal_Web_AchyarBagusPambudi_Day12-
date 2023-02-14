let localStorageKey = "KEY";

function checkStorage() {
  return typeof Storage !== undefined;
}

function formValidation(nama, email, numberPhone, subject, textArea) {
  if (nama == "") {
    alert("harap masukan nama");
  } else if (email == "") {
    alert("harap masukan email");
  } else if (numberPhone == "") {
    alert("harap masukan nomor");
  } else if (subject == "") {
    alert("harap masukan subject anda");
  } else if (textArea == "") {
    alert("harap masukan pesan anda");
  } else {
    let myEmail = "achyarbagus@gmail.com";
    let sendEmail = document.createElement("a");
    sendEmail.href = `mailto:${myEmail}? subject = ${subject} &body=hallo nama saya ${nama} dan ${textArea}`;
    sendEmail.click();
  }
}

function submitButton() {
  //
  let nama = document.getElementById("nama").value;
  let email = document.getElementById("email").value;
  let numberPhone = document.getElementById("phone").value;
  let subject = document.getElementById("subject").value;
  let textArea = document.getElementById("text-area").value;

  console.log(email);
  formValidation(nama, email, numberPhone, subject, textArea);

  //
}
