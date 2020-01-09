const signUpButton = document.getElementById('signUp');
const signInButton = document.getElementById('signIn');
const container = document.getElementById('container');
const linklogIn = document.getElementById('loginLink');
const linkSignUp = document.getElementById('signupLink');
 a = false;





function controller(x){
    
    var container = document.getElementById('container');

    if (x.matches) {

    	container.classList.remove("my-overlay-right");
        document.getElementById("s").innerHTML="king seccussful!!!";

         a = false; 
         container.classList.add("my-overlay-right");  

         linklogIn.addEventListener('click',() =>{
         	$("#signinM").show();
         	$("#signupM").hide();

         });


         linkSignUp.addEventListener('click',() =>{
         	$("#signinM").hide();
         	$("#signupM").show();
         });

        $("#signin").hide();
        $("#signup").hide();
        $("#signupM").show();
        $("#over").hide();


       
		
    }
   
    else{

    	$("#signinM").hide();

    	$("#signin").show();

        $("#over").show();
    	document.getElementById("s").innerHTML="king !!!";
    	   	 
    	 a = true;
    	
        signUpButton.addEventListener('click', () => {
			if (a) {
				container.classList.add("right-panel-active");
			}
			
		});

		signInButton.addEventListener('click', () => {

			if (a) {
				container.classList.remove("right-panel-active");
			}
			
		});		
    }
}


var b = window.matchMedia("(max-width:500px)")
controller(b)
b.addListener(controller)


