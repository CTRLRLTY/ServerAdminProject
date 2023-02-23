(function() {
    document.getElementById("email").innerText = "welcome " + get_cookie("session");

    fetch("/purchase-item-count")
        .then(async (response) => {
            switch (response.status)
            {
            case 200:
                let data = await response.json();

                for (let i = 0; i < data.length; ++i) 
                {
                    let columns = data[i];
                    let index = columns[0];
                    let amount = columns[1];

                    elemLabelCount = document.querySelector(`.item:nth-child(${index+1}) .item-count`);
                    elemLabelCount.innerText = amount;
                }

                break
            }
        });
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

async function handle_purchase(index, amount)
{
    let fdata = new FormData();
    fdata.append("index", index);
    fdata.append("email", get_cookie("session"));
    fdata.append("amount", amount);

    let res = await fetch("/purchase-item", {
        method: 'POST',
        body: fdata
    });

    if(res.ok) 
        location.reload();
}