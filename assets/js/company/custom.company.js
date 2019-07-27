$(document).ready(function(){
    $("#btn-create-company").click(function(){
        $("#modal-create-edit-company-title").html("Create new company");
        document.getElementById("create-edit-company-form").action = $(this).data('url');

        $("#modal-create-edit-company").modal("show");
    });

    $("#create-edit-company-form").validate({
        rules: {
            name: {
                required: true
            },
            slug: {
                required: true
            }
        },
        messages: {
            name: {
                required: "Name is a mandatory field."
            },
            slug: {
                required: "Slug is a mandatory field."
            }
        },
        submitHandler: function(form) {
            form.submit();    
            $("#modal-create-edit-company").modal("hide");        
            toggleLoading();
        }
    });
});
