var isMobile = navigator.userAgent.match(/Android|BlackBerry|iPhone|iPad|iPod|Opera Mini|IEMobile/i);
var drawer, origin, content;
function showDrawer(id, name, orgId) {
	origin.attr("href", id);
	origin.attr("data-origin", orgId);
	origin.text(name);
	content.html($(id + " > .reference-text").html());
	drawer.addClass("visible");
}

function getSettings(key, type) {
	if ("localStorage" in window) {
		if (type === "boolean") {
			return (window.localStorage[key] === "true");
		} else {
			return window.localStorage[key];
		}
	} else {
		var value = "; " + document.cookie;
		var parts = value.split("; " + key + "=");
		if (parts.length == 2) {
			if (type === "boolean") {
				return (parts.pop().split(";").shift() === "true");
			} else {
				return parts.pop().split(";").shift();
			}
		}
	}
}

function setSettings(key, value) {
	if ("localStorage" in window) {
		window.localStorage[key] = value;
	} else {
		var now = new Date();
		var time = now.getTime();
		var expireTime = time + (10 * 365 * 24 * 24);
		now.setTime(expireTime);
		document.cookie = key + "=" + value + ";expires=" + now.toGMTString() +";path=/";
	}
}

function initSettings() {
	if (document.getElementById("reference-settings")) {
		$("#reference-settings").show();
		return;
	}
	var setbox = $("<form>").attr("id", "reference-settings").append(
		$("<div>").addClass("content").append(
			$("<label>").append(
				$("<input>").attr({"type": "checkbox", "name": "showRefOnHover", "checked": getSettings("showRefOnHover", "boolean")}),
				"No Click footnote."
			)
		),
		$("<div>").addClass("foot").append(
			$("<input>").attr({"type": "button", "value": "Cancel"}).click(function() {
				$("#reference-settings").hide();
			}),
			$("<input>").addClass("save-settings").attr({"type": "submit", "value": "Save"})
		)
	).submit(function(e) {
		e.preventDefault();
		setSettings("showRefOnHover", $(this).find("[name=showRefOnHover]").is(":checked"));
		location.reload();
	});
	$(document.body).append(setbox);
}

function showTooltip(elem) {
	showDrawer(elem.parentNode.getAttribute("href"), elem.textContent, elem.parentNode.parentNode.getAttribute("id"));
	drawer.show();
	drawer.css({ top: ($(elem).offset().top - drawer.outerHeight()), left: $(elem).offset().left });
	drawer.stop().animate({opacity: 1}, 100);
}

function hideTooltip() {
	drawer.removeClass("visible");
	drawer.animate({opacity: 0}, 100, function() { $(this).hide(); });
}

$(document).ready(function($) {
	if (getSettings("showRefOnHover") === undefined) {
		setSettings("showRefOnHover", true);
	}
	if (isMobile) {
		$(document.body).addClass("mode-drawer");
	} else {
		$(document.body).addClass("mode-tooltip");
	}
	/* create drawer */
	drawer = $("<div>").attr("id", "reference-drawer");
	origin = $("<a>").attr("id", "reference-origin").click(function() {
		$('html, body').animate({scrollTop: ($($(this).attr("href")).offset().top - 60)}, 400);
	});
	content = $("<span>").attr("id", "reference-drawer-text");
	var settingsIcon = $("<span>").addClass("settings-icon").click(function() {
		initSettings();
	});
	var closeDrawer = $("<span>").addClass("close-icon").click(function() {
		drawer.removeClass("visible");
	});
	drawer.append(settingsIcon);
	drawer.append(closeDrawer);
	drawer.append(origin);
	drawer.append(content);
	$(document.body).append(drawer);
	
	$(document).click(function(e) {
		if (!$(e.target).closest("#reference-drawer").length)  {
			if (!isMobile) {
				hideTooltip();
			} else {
				drawer.removeClass("visible");
			}
		}
	}).scroll(function(e) {
		if ("ontouchstart" in window && isMobile) {
			drawer.removeClass("visible");
		}
	});
	$(".reference a").each(function() {
		var span = document.createElement("span");
		span.className = "reference-hooker";
		span.appendChild(this.childNodes[0]);
		this.appendChild(span);
	});
	if (getSettings("showRefOnHover", "boolean") && !isMobile) {
		drawer.hover(function(e) {
			showTooltip($("#" + origin.attr("data-origin") + " .reference-hooker").get(0));
		}, function(e) {
			hideTooltip();
		});
		$(".reference-hooker").hover(function(e) {
			showTooltip(this);
		}, function(e) {
			hideTooltip();
		});
	}
	$(".reference-hooker").click(function(e) {
		if (isMobile) {
			e.preventDefault();
			e.stopPropagation();
			showDrawer(this.parentNode.getAttribute("href"), this.textContent);
			drawer.addClass("visible");
		} else if (!getSettings("showRefOnHover", "boolean")) {
			e.preventDefault();
			e.stopPropagation();
			showTooltip(this);
		}
	});
});
/* 오리위키에서 가져옴 */