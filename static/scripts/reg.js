const regBtn = document.querySelector('.reg-btn')

regBtn.addEventListener('click', async(regEvent)=>{
    regEvent.preventDefault()
    const login = document.querySelector('[name="email"]').value
    const password = document.querySelector('[name="password"]').value
    const passwordRepeat = document.querySelector('[name="password-repeat"]').value
    console.log(login, password, passwordRepeat)
})