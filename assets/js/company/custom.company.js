$(document).ready(function(){    
    var $form = document.getElementById("create-edit-company-form");
    var $deleteform = document.getElementById("delete-company-form");
    var $helpBlock = document.getElementById("slug-help-block");

    // Create / edit the company
    $("#btn-create-company").click(function(){
        $("#modal-create-edit-company-title").html("Create new company");
        $form.action = $(this).data('url');
        $form.reset();
        $helpBlock.innerHTML = "";

        $("#modal-create-edit-company").modal("show");
    });

    $("#btn-edit-company").click(function(){
        $("#modal-create-edit-company-title").html("Edit company");
        $form.action = $(this).data('url');
        $helpBlock.innerHTML = "";

        $("#modal-create-edit-company").modal("show");
    });

    // Get the company detail when click on edit button
    $(".btn-edit-company").click(function(){
        // Show loading first
        toggleLoading();
        var id = $(this).data("id");
        var url = "/dashboard/company/" + id + "/show/json";
        var editURL = "/dashboard/company/" + id + "/update";
        var data;
        axios.get(url)
        .then(function (response) {
            // handle success
            $form.reset();
            $helpBlock.innerHTML = "";

            // Set the value in the dialog
            $form.action = editURL;
            data = response["data"];
            
            // Unhide the loading
            toggleLoading();   
            
            if(data["data"] != null) {
                $("#modal-create-edit-company-title").html(data["data"]["Name"]);
                $form.elements.namedItem("companyId").value = data["data"]["ID"];
                $form.elements.namedItem("name").value = data["data"]["Name"];
                $form.elements.namedItem("slug").value = data["data"]["Slug"];
                $form.elements.namedItem("description").value = data["data"]["Description"];
                $form.elements.namedItem("email").value = data["data"]["Email"];
                $form.elements.namedItem("phone").value = data["data"]["Phone"];
                $form.elements.namedItem("fax").value = data["data"]["Fax"];
                $form.elements.namedItem("address").value = data["data"]["Address"];
                $("#modal-create-edit-company").modal("show");
            } else {
                Swal.fire({
                    type: 'error',
                    title: 'Oops...',
                    text: 'Something went wrong! Please refresh the page.',
                })
            }
                     
        })
        .catch(function (error) {
            // handle error
            toggleLoading();
            Swal.fire({
                type: 'error',
                title: 'Oops...',
                text: 'Something went wrong!',
            })
        });
    });  

    $.validator.addMethod('url', function (value) { 
        return /^[a-zA-Z0-9_]*$/.test(value); 
    }, 'Only alphabets and numerics are allowed.');  

    var validator = $("#create-edit-company-form").validate({
        rules: {
            name: {
                required: true
            },
            slug: {
                required: true,
                url: true
            }
        },
        messages: {
            name: {
                required: "Name is a mandatory field."
            },
            slug: {
                required: "URL is a mandatory field.",
                url: "Only alphabets and numerics are allowed."
            }
        },
        submitHandler: function(form) {
            form.submit();    
            $("#modal-create-edit-company").modal("hide");        
            toggleLoading();
        }
    });
    
    // Set the URL for delete when click on delete button
    $(".btn-delete-company").click(function(){        
        var id = $(this).data("id");
        var deleteURL = "/dashboard/company/" + id + "/delete";
        deleteCompany(deleteURL);
    });

    $("#btn-delete-company").click(function(){
        var deleteURL  = $(this).data('url');
        deleteCompany(deleteURL);
    });

    // Get unique URL
    // Get the input box
    var slugInput = document.getElementById('inputSlug');

    // Init a timeout variable to be used below
    var timeout = null;

    // Listen for keystroke events
    slugInput.onkeyup = function (e) {

        // Clear the timeout if it has already been set.
        // This will prevent the previous task from executing
        // if it has been less than <MILLISECONDS>
        clearTimeout(timeout);

        // Make a new timeout set to go off in 800ms
        timeout = setTimeout(function () {
            if(validator.check('#inputSlug')){
                /*field is valid*/
                $helpBlock.innerHTML = '';
                // Get the unique URL
                var compId = $form.elements.namedItem("companyId").value;
                var slug = slugInput.value;
                var url = "/dashboard/company/getUniqueSlug?comp=" + compId + "&slug=" + slug;
                axios.get(url)
                .then(function (response) {
                    // handle success 
                    data = response["data"];
                    
                    if(!data["is_unique"]) {
                        $helpBlock.innerHTML = '<i class="text-danger fa fa-exclamation-circle"></i> The URL has already been taken.';
                    } else {
                        $helpBlock.innerHTML = '<i class="text-success fa fa-check-circle"></i> The URL is available.';
                    }
                })
                .catch(function (error) {
                    // handle error
                    console.log(error);
                });
            } else {
                /*field is not valid (but no errors will be displayed)*/
                $helpBlock.innerHTML = '<i class="text-danger fa fa-exclamation-circle"></i> Invalid URL.';
            }

        }, 500);
    };
    
    function deleteCompany(url) {
        $deleteform.action = url;
        Swal.fire({
            title: 'Are you sure?',
            text: "All the relevant data will be deleted!",
            type: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Yes, delete it!'
        }).then((result) => {
            if (result.value) {   
                $deleteform.submit();
                toggleLoading();
            }
        })
    }
});
