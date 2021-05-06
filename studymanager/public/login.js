  function userRegister() {
    let nameRegisterInput = $('#nameRegister').val();
    let emailRegisterInput = $('#emailRegister').val();
    let passwordRegisterInput = $('#passwordRegister').val();

    $.ajax({
        type: 'POST',
        url: '/users',
        dataType: 'html',
        data : {'name': nameRegisterInput, 'email':emailRegisterInput, 'pass_word':passwordRegisterInput},
        success: function (data) {
            alert('회원가입이 완료 되었습니다.');
        }
    });

}

 function userLogin() {
    let emailLoginInput = $('#emailLogin').val();
    let passwordLoginInput = $('#passwordLogin').val();

    $.ajax({
        type: 'POST',
        url: '/login',
        dataType: 'html',
        async: false,
        data: {'email':emailLoginInput, 'pass_word':passwordLoginInput},
        success: function (data) {
            console.log("todos로 이동");
            //window.location.href = "/todo";
        }
    });

}