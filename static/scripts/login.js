const loginBtn = document.querySelector('button')

loginBtn.addEventListener('click', (event)=>{
    event.preventDefault()
    const email = document.querySelector('[name="email"]').value
    const password = document.querySelector('[name="password"]').value
    console.log(email, password)
})