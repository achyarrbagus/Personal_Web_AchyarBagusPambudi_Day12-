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
    let emailDefault = "uzumaki-akbar@gmail.com";
    let sendEmail = document.createElement("a");
    sendEmail.href = `mailto:${emailDefault}? subject = ${subject} &body=hallo nama saya ${nama} dan ${textArea}`;
    sendEmail.click();
  }
}

function contactMe(event) {
  event.prevenDefault();
  let yourName = document.getElementById("yourName").value;
  let yourEmail = document.getElementById("yourEmail").value;
  let yourNumber = document.getElementById("yourNumber").value;
  let selectCarrer = document.getElementById("selectCarrer").value;
  let descriptionMail = document.getElementById("descriptionMail").value;

  formValidation(yourName, yourEmail, yourNumber, selectCarrer, descriptionMail);
}
