var privateIpfsGateway = "http://127.0.0.1:8080/ipfs/";
var ipfsGateway = "";

window.addEventListener('DOMContentLoaded', (event) => {
	ipfsGateway = localStorage.getItem('ipfsGateway');

	if (ipfsGateway) {
		document.querySelectorAll("div.post img, audio, video").forEach(function(e){
			var cid = e.src.substring(e.src.length - 46);
			e.src = ipfsGateway + cid;
		});
	}

	document.querySelectorAll("span.upvote a, span.downvote a").forEach(function(e){
		e.onclick = vote;
	});
});

document.addEventListener('keydown', (event) => {
	if (event.ctrlKey && event.key === 'i') {
		var ipfsGateway = window.prompt("Enter new IPFS gateway URL (or empty string to use default)", privateIpfsGateway);
		if (ipfsGateway !== null) {
			localStorage.setItem('ipfsGateway', ipfsGateway);
			window.location.reload();
		}
	}
});

function vote() {
	var newHTML = this.textContent + " " + (Number(this.parentNode.lastChild.textContent.trim()) + 1);
	var parentNode = this.parentNode;

	var xmlHTTP = new XMLHttpRequest();
	xmlHTTP.addEventListener('load', function() {
		if (this.status == 201) parentNode.innerHTML = newHTML;
	});
	xmlHTTP.open("GET", this.href);
	xmlHTTP.send();

	return false;
}

