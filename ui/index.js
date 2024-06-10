let loginForm = document.getElementById("upload-form");
loginForm.addEventListener("submit", (e) => {
    e.stopPropagation()
    e.preventDefault()


    const file = document.getElementById("file").files[0]

    const reader = new FileReader()
    reader.onload = (e) => {
        let data = e.target.result
        let workbook = XLSX.read(data)

        const worksheet = workbook.Sheets[workbook.SheetNames[0]];
        console.log(worksheet)

        // const table = XLSX.utils.sheet_to_html(worksheet)
        // console.log(table)

        const json = XLSX.utils.sheet_to_json(worksheet, {header: 1});
        console.log(json)

        const headers = json[0]
        console.log(headers)

        const headersEl = document.getElementById("headers")
        for (const header of headers) {

            headersEl.appendChild(document.createElement(`<td>${header}</td>`))
        }
    }

    reader.readAsArrayBuffer(file)

})
