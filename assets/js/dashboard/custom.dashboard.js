$(document).ready(function(){
    
    // Get the posts    
    const postLimitPerLoad = 10;
    var lastPublishedID = "";
    var lastPublishedDT = "";
    var isPublishedPostLoading = false;
    var isLastPublishedPostShown = false;
    var $publishedPostsContainer = $('#dashboard-posts');  

    function loadPublishedPostAjax(url) {        
        // Set the loading icon...  
        if(!isPublishedPostLoading 
            && !isLastPublishedPostShown
        ) {
            isPublishedPostLoading = true;
            $publishedPostsContainer.append( '<p class="mt-2 text-center text-muted post-loading-container"><i class="fa fa-spinner"></i> Loading...</p>' );
            
            axios.get(url)
            .then(function (response) {
                // handle success    
                console.log(response);
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

    // Get the published post on initial load
    if($publishedPostsContainer.length > 0) {
        loadPost();
    } 

    // Detecting scroll to bottom of div
    $(window).bind('scroll', function() {
        loadPost();
    });
    
    // Load the post
    function loadPost() {
        var additionalFilter = "";
        // Dashboard post
        if( $(window).scrollTop() >= $publishedPostsContainer.offset().top + $publishedPostsContainer.outerHeight() - window.innerHeight ) {
            
            if(lastPublishedID != "" && lastPublishedDT != "") {
                additionalFilter = "&lastID=" + lastPublishedID + "&lastPublished=" + lastPublishedDT;
            }

            loadPublishedPostAjax("/dashboard/post/list?status=2&limit=" + postLimitPerLoad + additionalFilter);
        }
    }
});