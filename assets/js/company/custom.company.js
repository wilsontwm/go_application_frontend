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

    // Init a timeout variable to be used below
    var slugTimeout = null;
    var nameTimeout = null;
    // Get the URL based on the input name
    $('#inputName').keyup(function () {

        // Clear the timeout if it has already been set.
        // This will prevent the previous task from executing
        // if it has been less than <MILLISECONDS>
        clearTimeout(nameTimeout);

        // Make a new timeout set to go off in 800ms
        nameTimeout = setTimeout(function () {
            var value = $("#inputName").val();
            var slugValue = value.trim().replace(/\s+/g, '-').toLowerCase();
            $("#inputSlug").val(slugValue);
            // Get the slug availability
            getSlugAvailability(slugTimeout);
        }, 500);
    });

    $.validator.addMethod('url', function (value) { 
        return /^[a-zA-Z0-9-]*$/.test(value); 
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

    
    // Listen for keystroke events
    slugInput.onkeyup = function (e) {
        getSlugAvailability(slugTimeout);
    };

    function getSlugAvailability(timeout) {
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
    }
    
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

    // Invite users to the company
    var $inviteForm = document.getElementById("invite-to-company-form");
    $("#btn-invite-to-company").click(function(){
        $inviteForm.reset();

        $("#modal-invite-to-company").modal("show");
    });
    
    var $inputInvitationEmail = document.getElementById("inputInvitationEmail");
    if (typeof($inputInvitationEmail) != 'undefined' && $inputInvitationEmail != null) {
        $("#inputInvitationEmail").multiple_emails();

        $('#inputInvitationEmail').change(function(){
            $("#inputEmails").val($(this).val());
        });
    }    

    $("#btn-submit-invite-to-company").click(function(e){
        e.preventDefault();
        
        // Check the size of the array        
        if($("#inputEmails").val() != "" &&  JSON.parse($("#inputEmails").val()).length > 0) {
            $inviteForm.submit();
            toggleLoading();            
        } else {
            Swal.fire({
                type: 'error',
                title: 'Oops...',
                text: 'You have not keyed in any email yet.',
            });
        }

        $("#modal-invite-to-company").modal("hide"); 
    });

    // Get the invitation result via AJAX    
    function invitationTemplate(data) {
        var html = '';
        $.each(data, function(index, item){
            var label = '<small class="label label-primary">Awaiting Response</small>';
            var disabled = '';
            if(item["Status"] == 1) {
                label = '<small class="label label-success">Accepted</small>';
                disabled = 'disabled';
            } else if(item["Status"] == 2) {
                label = '<small class="label label-danger">Declined</small>';
                disabled = 'disabled';
            }

            html += '<tr>'
                    + '<td>'+ item["Email"] +'</td>'
                    + '<td>' + label + '</td>'
                    + '<td><button class="btn btn-default btn-resend-invitation mr-1" data-id="'+ item["ID"] +'" data-company-id="'+ item["CompanyID"] +'" ' + disabled + '>Resend invitation</button>'
                    + '<button class="btn btn-danger btn-delete-invitation" data-id="'+ item["ID"] +'" data-company-id="'+ item["CompanyID"] +'">Delete</button></td>'
                    + '</tr>';
        });

        return html;
    }

    function loadInvitationAjax(url) {
        var invitationURL = url;
        axios.get(invitationURL)
        .then(function (response) {
            // handle success             
            data = response["data"];
            $('#invitation-pagination-container').pagination({
                pageSize: 25,
                showGoInput: true,
                showGoButton: true,
                dataSource: data["data"],
                callback: function(d, pagination) {
                    // template method of yourself
                    var html = invitationTemplate(d);
                    $('#invitation-results-container').html(html);
                }
            });  
        })
        .catch(function (error) {
            // handle error
            console.log(error);
        });    
    }

    var $invitationpane = document.getElementById("invitation-pane");
    if (typeof($invitationpane) != 'undefined' && $invitationpane != null) {
        var url = $invitationpane.getAttribute("data-url");
        loadInvitationAjax(url);
    }

    // Resend the invitation email to the invites
    // Init a timeout variable to be used below
    var resendtimeout = null;
        
    $(document).on('click', '.btn-resend-invitation', function(){ 
        // Clear the timeout if it has already been set.
        // This will prevent the previous task from executing
        // if it has been less than <MILLISECONDS>
        clearTimeout(resendtimeout);

        var id = $(this).data("id");
        var compId = $(this).data("company-id");
        let csrfToken = document.getElementsByName("gorilla.csrf.Token")[0].value;
        // Make a new timeout set to go off in 800ms
        resendtimeout = setTimeout(function () {
            toggleLoading();
            var inviteURL = "/dashboard/company/" + compId + "/invite/" + id + "/resend";
            
            axios({
                method: 'post',
                url: inviteURL,
                headers: {"X-CSRF-Token": csrfToken},
            })
            .then(function (response) {
                // handle success
                data = response["data"];
                
                // Unhide the loading
                toggleLoading();   
                if(data["success"]) {
                    Swal.fire({
                        type: 'success',
                        title: 'Awesome!',
                        text: 'You have successfully resend invitation email to ' + data["data"]["Email"] + '.',
                    })
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
                console.log(error);
                Swal.fire({
                    type: 'error',
                    title: 'Oops...',
                    text: 'Something went wrong!',
                })
            });
        }, 500);
    });

    // Delete the company invitation
    // Set the URL for delete when click on delete button
    $(document).on('click', '.btn-delete-invitation', function(){        
        var id = $(this).data("id");        
        var compId = $(this).data("company-id");
        var deleteInvitationURL = "/dashboard/company/" + compId + "/invite/" + id + "/delete";
        var $row = $(this).parents("tr");
        let csrfToken = document.getElementsByName("gorilla.csrf.Token")[0].value;
        
        Swal.fire({
            title: 'Are you sure?',
            text: "The invitation will be deleted!",
            type: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Yes, delete it!'
        }).then((result) => {
            if (result.value) {   
                toggleLoading();

                axios({
                    method: 'post',
                    url: deleteInvitationURL,
                    headers: {"X-CSRF-Token": csrfToken},
                })
                .then(function (response) {
                    // handle success                    
                    // Unhide the loading
                    toggleLoading();   
                    if(response["status"] == 200) {
                        $row.remove();
                        Swal.fire({
                            type: 'success',
                            title: 'Awesome!',
                            text: 'You have successfully removed the invitation.',
                        })
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
                    console.log(error);
                    Swal.fire({
                        type: 'error',
                        title: 'Oops...',
                        text: 'Something went wrong!',
                    })
                });
            }
        })
        
    });
});
