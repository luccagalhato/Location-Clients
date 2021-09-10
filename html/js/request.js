async function getClientsNear(body, path) {
    var myInit = {
      method: "POST",
      headers: new Headers(),
      mode: "cors",
      cache: "default",
      body: body,
    };
    const response = await fetch(`${window.location.protocol}//${window.location.host}/request`, myInit);
    const blob = await response.blob();
    const text = await blob.text();
    const clientNear = JSON.parse(text)
    console.log(clientNear)
  }
//   function enviarData() {
//     var input = document.getElementById("dateInicial");
//     var input2 = document.getElementById("dateFinal");
//     var date = {
//       dataInicial: input.value,
//       dataFinal: input2.value,
//     };
//     var datas = JSON.stringify(date);
//     getIds(datas, "date");
//   }
  
//   function enviarGtin() {
//     var input = document.getElementById("text");
//     getIds(input.value, "gtin");
//   }
  
  