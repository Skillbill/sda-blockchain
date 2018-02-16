const qs = document.querySelector.bind(document);

function dapp_init() {
    const el = qs('#providerAddr');
    const addr = el.value;
    if (!addr) {
	el.classList.add('error');
	return;
    }
    el.classList.remove('error');
    init(addr);
    prepareDeployBox();
}

function dapp_deploy() {
    let el = qs('#contractOwner');
    const owner = el.value;
    if (!Web3.utils.isAddress(owner)) {
	el.classList.add('error');
	return;
    }
    el.classList.remove('error');
    const candidates = getCandidates({validate: true});
    if (!candidates) {
        return;
    }
    deploy(owner, candidates).then(txaddr => {
	prepareVotingBox(txaddr, candidates);
    }).catch (err => alert(err));
}

function dapp_voteFor(n) {
    let el = qs('#voterAddress');
    const addr = el.value;
    if (!Web3.utils.isAddress(addr)) {
	el.classList.add('error');
	return;
    }
    el.classList.remove('error');
    const candidate = getCandidate(n);
    voteFor(candidate, addr).then(() => {
	el.value = '';
    }).catch(err => alert(err));
}

function dapp_results() {
    hide('#resultButton');
    hide('#voter');
    hide('.vote');

    const candidates = getCandidates()
    for (let i=0; i<candidates.length; i++) {
	totalVotesFor(candidates[i]).then(count => {
	    let td = qs(`#votecell_${i}`);
	    td.innerText = count;
	}).catch(err => alert(err));
    }
}

function getCandidates(options) {
    const check = options && options.validate;
    const candidates = document.getElementsByClassName('candidateBox');
    let lst = [];
    for (i=0; i<candidates.length; i++) {
	let child = candidates[i].childNodes[0];
	if (check && !child.value) {
	    child.classList.add('error');
	    return null;
	}
        child.classList.remove('error');
	lst.push(child.value);
    }
    return lst;
}

function getCandidate(n) {
    const el = qs(`#cell_candidate_${n}`);
    return el.innerText;
}

function addCandidate() {
    const els = document.getElementsByClassName('candidateBox');
    const box = buildCandidateBox(els.length);
    const parent = qs('#candidates');
    parent.append(box);
}

function removeCandidate(n) {
    if (n < 2) {
        alert('FIXME, n<2 in removeCandidate()');
        return;
    }
    let el = qs(`#candidate_${n}`);
    el.remove();
}

function buildCandidateBox(n) {
    let div = document.createElement('DIV');
    div.classList.add('candidateBox');
    div.id = `candidate_${n}`;

    let input = document.createElement('INPUT');
    input.classList.add('candidate');
    input.type = 'text';
    div.append(input);

    if (n > 1) {
        let remove = document.createElement('INPUT');
        remove.type = 'button';
        remove.value = '-';
        remove.addEventListener('click', function() { removeCandidate(n); });
        div.append(remove);
    }
    return div;
}

function prepareDeployBox() {
    hide('#connectBox');
    unhide('#deployBox');
    const el = qs('#candidates');
    addCandidate();
    addCandidate();
}

function createVoteButton(n) {
    let el = document.createElement('INPUT');
    el.classList.add('vote');
    el.id = `submit_vote_${n}`;
    el.type = 'button';
    el.value = 'vota';
    el.addEventListener('click', function() { dapp_voteFor(n); });
    return el;
}

function prepareVotingBox(contractAddr, candidates) {
    hide('#deployBox');
    unhide('#votingBox');
    qs('#votingTitle').innerText = `Contract: ${contractAddr}`;

    const table = qs('#votingTable');
    for (let i=0; i<candidates.length; i++) {
	let row = table.insertRow(i+1);
	let nameCell = row.insertCell(0);
	nameCell.classList.add('candidate');
	nameCell.innerText = candidates[i];
	nameCell.id = `cell_candidate_${i}`;
	let voteForCell = row.insertCell(1);
	voteForCell.id= `votecell_${i}`;
	voteForCell.append(createVoteButton(i));
    }
}

function hide(selector) {
    const v = document.querySelectorAll(selector);
    v.forEach(el => el.classList.add('hidden'));
}

function unhide(selector) {
    const v = document.querySelectorAll(selector);
    v.forEach(el => el.classList.remove('hidden'));
}
