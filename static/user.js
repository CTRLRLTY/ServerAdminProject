function handle_user(ev, endpoint)
{
    ev.preventDefault();

    const form = new FormData(document.forms[0]);

    fetch(endpoint, {method: 'POST', body: form})
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