$(document).ready(function(){
    var userId = getUserIdFromURL();

    // Datepicker
    $("#inputBirthday").datetimepicker({
        format: 'DD MMM YYYY',
    });

    $.validator.addMethod("password",function(value,element){
        return this.optional(element) || /^(?=.*\d)(?=.*[a-z])(?=.*[A-Z]).{8,16}$/i.test(value);
    },"Passwords are 8-16 characters with uppercase letters, lowercase letters and at least one number.");
    
    $("#edit-profile-form").validate({
        rules: {
            name: {
                required: true
            },
            birthday: {
                date: true
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
        if (this.files && this.files[0]) {
            var reader = new FileReader();
            reader.readAsDataURL(this.files[0]);

            $("#upload-profile-picture-form").submit();            
        }
    });

    $("#upload-profile-picture-form").on('submit', function(){
        toggleLoading();
    });
    
    $("#btn-delete-profile-picture").click(function(e){
        e.preventDefault();
        Swal.fire({
            title: 'Are you sure?',
            text: "You will not be able to recover the picture once deleted!",
            type: 'warning',
            showCancelButton: true,
            confirmButtonColor: '#3085d6',
            cancelButtonColor: '#d33',
            confirmButtonText: 'Yes, delete it!'
        }).then((result) => {
            if (result.value) {   
                $("#delete-profile-picture-form").submit();
                toggleLoading();
            }
        })
    });

    // Get the posts    
    const postLimitPerLoad = 10;
    var lastPublishedID = "";
    var lastPublishedDT = "";
    var isPublishedPostLoading = false;
    var isLastPublishedPostShown = false;
    var $publishedPostsContainer = $('#published-posts');        
    var lastDraftID = "";
    var lastDraftDT = "";
    var isDraftPostLoading = false;
    var isLastDraftPostShown = false;
    var $draftPostsContainer = $('#draft-posts');
    var lastScheduledID = "";
    var lastScheduledDT = "";
    var isScheduledPostLoading = false;
    var isLastScheduledPostShown = false;
    var $scheduledPostsContainer = $('#scheduled-posts');

    function loadPublishedPostAjax(url) {        
        // Set the loading icon...  
        if(!isPublishedPostLoading 
            && !isLastPublishedPostShown
            && $publishedPostsContainer.hasClass("active")
        ) {
            isPublishedPostLoading = true;
            $publishedPostsContainer.append( '<p class="mt-2 text-center text-muted post-loading-container"><i class="fa fa-spinner"></i> Loading...</p>' );
            
            axios.get(url)
            .then(function (response) {
                // handle success    
                //console.log(response);
                data = response["data"]["data"];
    
                // Remove the loading container
                $publishedPostsContainer.children(".post-loading-container").remove();
                
                if(data != null) {
                    for(i = 0; i < data.length; i++) {
                        var ele = data[i];
                        var html = postTemplate(ele, response["data"]["defaultProfilePic"]);
                        $publishedPostsContainer.append(html);
        
                        // Set the last viewed post ID and published date
                        lastPublishedID = ele.ID;
                        lastPublishedDT = ele.PublishedAt;
                    }
                }   

                if(data == null || data.length == 0) { 
                    isLastPublishedPostShown = true; 
                }
                
                isPublishedPostLoading = false;

                // Load again if the bottom is still in view and has data
                loadPost();
            })
            .catch(function (error) {
                // handle error
                console.log(error);
            });   
        }    	
    }

    function loadDraftPostAjax(url) {        
        // Set the loading icon...  
        if(!isDraftPostLoading 
            && !isLastDraftPostShown
            && $draftPostsContainer.hasClass("active")
        ) {
            isDraftPostLoading = true;
            $draftPostsContainer.append( '<p class="mt-2 text-center text-muted post-loading-container"><i class="fa fa-spinner"></i> Loading...</p>' );
            
            axios.get(url)
            .then(function (response) {
                // handle success    
                //console.log(response);
                data = response["data"]["data"];
    
                // Remove the loading container
                $draftPostsContainer.children(".post-loading-container").remove();
                
                if(data != null) {
                    for(i = 0; i < data.length; i++) {
                        var ele = data[i];
                        var html = postTemplate(ele, response["data"]["defaultProfilePic"]);
                        $draftPostsContainer.append(html);
        
                        // Set the last viewed post ID and updated date
                        lastDraftID = ele.ID;
                        lastDraftDT = ele.UpdatedAt;
                    }
                }   

                if(data == null || data.length == 0) { 
                    isLastDraftPostShown = true; 
                }
                
                isDraftPostLoading = false;

                // Load again if the bottom is still in view and has data
                loadPost();
            })
            .catch(function (error) {
                // handle error
                console.log(error);
            });   
        }    	
    }

    function loadScheduledPostAjax(url) {        
        // Set the loading icon...  
        if(!isScheduledPostLoading 
            && !isLastScheduledPostShown
            && $scheduledPostsContainer.hasClass("active")
        ) {
            isScheduledPostLoading = true;
            $scheduledPostsContainer.append( '<p class="mt-2 text-center text-muted post-loading-container"><i class="fa fa-spinner"></i> Loading...</p>' );
            
            axios.get(url)
            .then(function (response) {
                // handle success    
                //console.log(response);
                data = response["data"]["data"];
    
                // Remove the loading container
                $scheduledPostsContainer.children(".post-loading-container").remove();
                
                if(data != null) {
                    for(i = 0; i < data.length; i++) {
                        var ele = data[i];
                        var html = postTemplate(ele, response["data"]["defaultProfilePic"]);
                        $scheduledPostsContainer.append(html);
        
                        // Set the last viewed post ID and updated date
                        lastScheduledID = ele.ID;
                        lastScheduledDT = ele.UpdatedAt;
                    }
                }   

                if(data == null || data.length == 0) { 
                    isLastScheduledPostShown = true; 
                }
                
                isScheduledPostLoading = false;

                // Load again if the bottom is still in view and has data
                loadPost();
            })
            .catch(function (error) {
                // handle error
                console.log(error);
            });   
        }    	
         
    }

    // Get the published post on initial load
    if($publishedPostsContainer.length > 0) {
        loadPost();
    } 

    // Detecting click on the nav bar and load accordingly
    $(document).on("click", ".nav-link", function() {         
        loadPost();
    }); 

    // Detecting scroll to bottom of div
    $(window).bind('scroll', function() {
        loadPost();
    });
    
    // Load the post
    function loadPost() {
        var additionalFilter = "";
        // Published post
        if( $publishedPostsContainer.hasClass("active") 
            && $(window).scrollTop() >= $publishedPostsContainer.offset().top + $publishedPostsContainer.outerHeight() - window.innerHeight ) {
            if(lastPublishedID != "" && lastPublishedDT != "") {
                additionalFilter = "&lastID=" + lastPublishedID + "&lastPublished=" + lastPublishedDT;
            }
            loadPublishedPostAjax("/dashboard/post/list?author=" + userId + "&status=2&limit=" + postLimitPerLoad + additionalFilter);
        }
        // Draft post
        else if( $draftPostsContainer.hasClass("active") 
                 && $(window).scrollTop() >= $draftPostsContainer.offset().top + $draftPostsContainer.outerHeight() - window.innerHeight ) {
            if(lastDraftID != "" && lastDraftDT != "") {
                additionalFilter = "&lastID=" + lastDraftID + "&lastUpdated=" + lastDraftDT;
            }
            loadDraftPostAjax("/dashboard/post/list?author=" + userId + "&status=0&limit=" + postLimitPerLoad + additionalFilter);
        }
        // Scheduled post        
        else if( $scheduledPostsContainer.hasClass("active") 
                 && $(window).scrollTop() >= $scheduledPostsContainer.offset().top + $scheduledPostsContainer.outerHeight() - window.innerHeight ) {
            if(lastScheduledID != "" && lastScheduledDT != "") {
                additionalFilter = "&lastID=" + lastScheduledID + "&lastUpdated=" + lastScheduledDT;
            }
            loadScheduledPostAjax("/dashboard/post/list?author=" + userId + "&status=1&limit=" + postLimitPerLoad + additionalFilter);
        }
    }
    
    // Get the user ID / author ID from the path
    function getUserIdFromURL() {
        var pathName = window.location.pathname;
        return pathName.split("/dashboard/user/")[1];
    }

});
