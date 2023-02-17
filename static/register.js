function handle_create_account(ev)
{
    ev.preventDefault();

    const form = new FormData(document.forms[0]);
    elemInfo.style.display = "block";

    fetch("/create-account", {method: 'POST', body: form})
        .then(async (response) => {
            switch (response.status)
            {
            case 200:
                let data = await response.json();
                elemInfo = document.getElementById("info");
                elemInfo.innerText = data.reason;
                elemInfo.style.display = "block";
                break
            case 303:
                window.location.assign(response.headers.get("Goto"));
                break
            }
        })
}