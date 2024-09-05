function UpdateMaskField(ele_id: string, value: string) {
  if (document.getElementById(ele_id) == null) {
    let ele = document.createElement("input")
    ele.hidden = true
    ele.id = ele_id
    ele.name = "UpdateMask"
    ele.value = value
    document.getElementById("UserUpdateForm")?.appendChild(ele)
  }
}
