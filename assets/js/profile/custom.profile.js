$(document).ready(function(){
    $.validator.addMethod("password",function(value,element){
        return this.optional(element) || /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,16}$/i.test(value);
    },"Passwords are 8-16 characters with uppercase letters, lowercase letters and at least one number.");
    
    $("#edit-profile-form").validate({
        rules: {
            name: {
                required: true
            }
        },
        messages: {
            name: {
                required: "Name is a mandatory field."
            }
        },
        submitHandler: function(form) {
            form.submit();            
            toggleLoading();
        }
    });

    $("#edit-password-form").validate({
        rules: {
            password: {
                required: true,
                password: true
            },
            retype_password: {
                equalTo: "#password_input"
            }
        },
        messages: {
            password: {
                required: "Password is a mandatory field.",
                password: "Passwords are 8-16 characters with uppercase letters, lowercase letters and at least one number."
            },
            retype_password: {
                equalTo: "Retype password does not match password."
            }
        },
        submitHandler: function(form) {
            form.submit();            
            toggleLoading();
        }
    });

    // Upload profile picture
    $("#btn-upload-picture").click(function() {
        $('#profile-image-input').click();
    });

    $("#profile-image-input").on('change',function(){
        console.log('Has profile image changed');
        if (this.files && this.files[0]) {
            console.log("Has file to submit");
            var reader = new FileReader();
            reader.readAsDataURL(this.files[0]);

            $("#upload-profile-picture-form").submit();            
        }
    });

    $("#upload-profile-picture-form").on('submit', function(){
        toggleLoading();
    });

});
