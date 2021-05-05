  function userRegister() {
    let nameRegisterInput = $('#nameRegister').val();
    let emailRegisterInput = $('#emailRegister').val();
    let passwordRegisterInput = $('#passwordRegister').val();

    let data = {'name': nameRegisterInput, 'email':emailRegisterInput, 'pass_word':passwordRegisterInput};

    $.ajax({
        type: 'POST',
        url: '/users',
        contentType: 'application/json',
        data: JSON.stringify(data),
        success: function (response) {
            alert('회원가입이 완료 되었습니다.');
        }
    });

}