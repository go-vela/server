// SPDX-License-Identifier: Apache-2.0

package constants

// Constants for build badges.
//
//nolint:godot // due to providing pretty printed svgs
const (
	// Badge for unknown state
	// <svg xmlns="http://www.w3.org/2000/svg" width="92" height="20">
	//     <linearGradient id="b" x2="0" y2="100%">
	//         <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
	//         <stop offset="1" stop-opacity=".1"/>
	//     </linearGradient>
	//     <path d="M0 3 a3 3 0 014-3h28v20H3 a3 3 0 01-3-3V3z" fill="#555555"/>
	//     <path d="M92 17 a3 3 0 01-3 3H32V0h56 a3 3 0 014 3v12z" fill="#9f9f9f"/>
	//     <rect width="100%" height="100%" rx="3" fill="url(#b)"/>
	//     <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
	//         <text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text>
	//         <text x="16" y="13" textLength="24" lengthAdjust="spacing">vela</text>
	//         <text x="62" y="14" fill="#010101" fill-opacity=".3" textLength="52" lengthAdjust="spacing">unknown</text>
	//         <text x="62" y="13" textLength="52" lengthAdjust="spacing">unknown</text>
	//     </g>
	// </svg>
	BadgeUnknown = `<svg xmlns="http://www.w3.org/2000/svg" width="92" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><path d="M0 3 a3 3 0 014-3h28v20H3 a3 3 0 01-3-3V3z" fill="#555555"/><path d="M92 17 a3 3 0 01-3 3H32V0h56 a3 3 0 014 3v12z" fill="#9f9f9f"/><rect width="100%" height="100%" rx="3" fill="url(#b)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text><text x="16" y="13" textLength="24" lengthAdjust="spacing">vela</text><text x="62" y="14" fill="#010101" fill-opacity=".3" textLength="52" lengthAdjust="spacing">unknown</text><text x="62" y="13" textLength="52" lengthAdjust="spacing">unknown</text></g></svg>`

	// Badge for success state
	// <svg xmlns="http://www.w3.org/2000/svg" width="85" height="20">
	//     <linearGradient id="a" x2="0" y2="100%">
	//         <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
	//         <stop offset="1" stop-opacity=".1"/>
	//     </linearGradient>
	//     <path d="M0 3 a3 3 0 014-3h30v20H3 a3 3 0 01-3-3V3z" fill="#555555"/>
	//     <path d="M85 17 a3 3 0 01-3 3H32V0h49 a3 3 0 014 3v12z" fill="#44cc11"/>
	//     <rect width="100%" height="100%" rx="3" fill="url(#a)"/>
	//     <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
	//         <text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text>
	//         <text x="16" y="13" textlength="24" lengthadjust="spacing">vela</text>
	//         <text x="58" y="14" fill="#010101" fill-opacity=".3" textlength="46" lengthadjust="spacing">success</text>
	//         <text x="58" y="13" textlength="46" lengthadjust="spacing">success</text>
	//     </g>
	// </svg>
	BadgeSuccess = `<svg xmlns="http://www.w3.org/2000/svg" width="85" height="20"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><path d="M0 3 a3 3 0 014-3h30v20H3 a3 3 0 01-3-3V3z" fill="#555555"/><path d="M85 17 a3 3 0 01-3 3H32V0h49 a3 3 0 014 3v12z" fill="#44cc11"/><rect width="100%" height="100%" rx="3" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text><text x="16" y="13" textlength="24" lengthadjust="spacing">vela</text><text x="58" y="14" fill="#010101" fill-opacity=".3" textlength="46" lengthadjust="spacing">success</text><text x="58" y="13" textlength="46" lengthadjust="spacing">success</text></g></svg>`

	// Badge for failed state
	// <svg xmlns="http://www.w3.org/2000/svg" width="73" height="20">
	//     <linearGradient id="a" x2="0" y2="100%">
	//         <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
	//         <stop offset="1" stop-opacity=".1"/>
	//     </linearGradient>
	//     <path d="M0 3 a3 3 0 014-3h30v20H3 a3 3 0 01-3-3V3z" fill="#555555"/>
	//     <path d="M73 17 a3 3 0 01-3 3H32V0h37 a3 3 0 014 3v12z" fill="#fe7d37"/>
	//     <rect width="100%" height="100%" rx="3" fill="url(#a)"/>
	//     <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
	//         <text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text>
	//         <text x="16" y="13" textlength="24" lengthadjust="spacing">vela</text>
	//         <text x="52" y="14" fill="#010101" fill-opacity=".3" textlength="46" lengthadjust="spacing">failed</text>
	//         <text x="52" y="13" textlength="46" lengthadjust="spacing">failed</text>
	//     </g>
	// </svg>
	BadgeFailed = `<svg xmlns="http://www.w3.org/2000/svg" width="73" height="20"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><path d="M0 3 a3 3 0 014-3h30v20H3 a3 3 0 01-3-3V3z" fill="#555555"/><path d="M73 17 a3 3 0 01-3 3H32V0h37 a3 3 0 014 3v12z" fill="#fe7d37"/><rect width="100%" height="100%" rx="3" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text><text x="16" y="13" textlength="24" lengthadjust="spacing">vela</text><text x="52" y="14" fill="#010101" fill-opacity=".3" textlength="46" lengthadjust="spacing">failed</text><text x="52" y="13" textlength="46" lengthadjust="spacing">failed</text></g></svg>`

	// Badge for error state
	// <svg xmlns="http://www.w3.org/2000/svg" width="69" height="20">
	//     <linearGradient id="a" x2="0" y2="100%">
	//         <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
	//         <stop offset="1" stop-opacity=".1"/>
	//     </linearGradient>
	//     <path d="M0 3 a3 3 0 014-3h30v20H3 a3 3 0 01-3-3V3z" fill="#555555"/>
	//     <path d="M69 17 a3 3 0 01-3 3H32V0h33 a3 3 0 014 3v12z" fill="#e05d44"/>
	//     <rect width="100%" height="100%" rx="3" fill="url(#a)"/>
	//     <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
	//         <text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text>
	//         <text x="16" y="13" textlength="24" lengthadjust="spacing">vela</text>
	//         <text x="50" y="14" fill="#010101" fill-opacity=".3" textlength="46" lengthadjust="spacing">error</text>
	//         <text x="50" y="13" textlength="46" lengthadjust="spacing">error</text>
	//     </g>
	// </svg>
	BadgeError = `<svg xmlns="http://www.w3.org/2000/svg" width="69" height="20"><linearGradient id="a" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><path d="M0 3 a3 3 0 014-3h30v20H3 a3 3 0 01-3-3V3z" fill="#555555"/><path d="M69 17 a3 3 0 01-3 3H32V0h33 a3 3 0 014 3v12z" fill="#e05d44"/><rect width="100%" height="100%" rx="3" fill="url(#a)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text><text x="16" y="13" textlength="24" lengthadjust="spacing">vela</text><text x="50" y="14" fill="#010101" fill-opacity=".3" textlength="46" lengthadjust="spacing">error</text><text x="50" y="13" textlength="46" lengthadjust="spacing">error</text></g></svg>`

	// Badge for running status
	// <svg xmlns="http://www.w3.org/2000/svg" width="88" height="20">
	//     <linearGradient id="b" x2="0" y2="100%">
	//         <stop offset="0" stop-color="#bbb" stop-opacity=".1"/>
	//         <stop offset="1" stop-opacity=".1"/>
	//     </linearGradient>
	//     <path d="M0 3 a3 3 0 014-3h28v20H3 a3 3 0 01-3-3V3z" fill="#555555"/>
	//     <path d="M88 17 a3 3 0 01-3 3H32V0h52 a3 3 0 014 3v12z" fill="#dfb317"/>
	//     <rect width="100%" height="100%" rx="3" fill="url(#b)"/>
	//     <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
	//         <text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text>
	//         <text x="16" y="13" textLength="24" lengthAdjust="spacing">vela</text>
	//         <text x="59" y="14" fill="#010101" fill-opacity=".3" textLength="46" lengthAdjust="spacing">running</text>
	//         <text x="59" y="13" textLength="46" lengthAdjust="spacing">running</text>
	//     </g>
	// </svg>
	BadgeRunning = `<svg xmlns="http://www.w3.org/2000/svg" width="88" height="20"><linearGradient id="b" x2="0" y2="100%"><stop offset="0" stop-color="#bbb" stop-opacity=".1"/><stop offset="1" stop-opacity=".1"/></linearGradient><path d="M0 3 a3 3 0 014-3h28v20H3 a3 3 0 01-3-3V3z" fill="#555555"/><path d="M88 17 a3 3 0 01-3 3H32V0h52 a3 3 0 014 3v12z" fill="#dfb317"/><rect width="100%" height="100%" rx="3" fill="url(#b)"/><g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11"><text x="16" y="14" fill="#010101" fill-opacity=".3" textLength="24" lengthAdjust="spacing">vela</text><text x="16" y="13" textLength="24" lengthAdjust="spacing">vela</text><text x="59" y="14" fill="#010101" fill-opacity=".3" textLength="46" lengthAdjust="spacing">running</text><text x="59" y="13" textLength="46" lengthAdjust="spacing">running</text></g></svg>`
)
