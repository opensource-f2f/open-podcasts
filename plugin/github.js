const meta = document.getElementsByClassName('gh-header-meta')[0]
const sourceBranch = meta.getElementsByTagName('a')[2]

const gitpod = document.createElement('a')
gitpod.href='https://gitpod.io/#' + sourceBranch.href
gitpod.innerText='Open with GitPod'
gitpod.innerHTML='<img src="https://gitpod.io/static/media/gitpod.2cdd910d.svg"/>'
gitpod.target='_blank'
meta.append(gitpod)
