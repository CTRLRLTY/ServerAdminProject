(function() {
    document.getElementById("email").innerText = "welcome " + get_cookie("session");
})();




function get_cookie(cname) {
  let name = cname + "=";
  let decodedCookie = decodeURIComponent(document.cookie);
  let ca = decodedCookie.split(';');

  for(let i = 0; i <ca.length; i++) {
    let c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }

  return "";
}

function handle_purchase(index, amount)
{
    let fdata = new FormData();
    fdata.append("index", index);
    fdata.append("email", get_cookie("session"));
    fdata.append("amount", amount);

    fetch("/purchase-item", {
        method: 'POST',
        // body: JSON.stringify({index, email: get_cookie("session"), amount})
        body: fdata
    })
        .then(async (response) => {
            switch (response.status)
            {
            case 200:
                let data = await response.json();
                elemLabelCount = document.querySelector(`.item-${data.index} item-count`);
                elemLabelCount.innerText = data.count;
                break
            }
        });
}