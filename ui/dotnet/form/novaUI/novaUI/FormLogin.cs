using System;
using System.Collections.Generic;
using System.ComponentModel;
using System.Data;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;
using System.Windows.Forms;

namespace novaUI
{
    public partial class FormLogin : Form
    {
        public FormLogin()
        {
            InitializeComponent();
        }

        private void buttonOK_Click(object sender, EventArgs e)
        {
            // verify user name validatation
            if (textBoxUserName.Text == "")
            {
                MessageBox.Show("用户名不能为空！");
                return;
            }
            if (textBoxUserName.TextLength < 8)
            {
                MessageBox.Show("用户名至少需要8位！");
                return;
            }
            // verify user password validation
            if (textBoxPassword.Text == "")
            {
                MessageBox.Show("密码不能为空！");
                return;
            }
            if (textBoxPassword.TextLength < 8)
            {
                MessageBox.Show("密码至少需要8位！");
                return;
            }
            if (!Regex.IsMatch(textBoxPassword.Text, @"^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*]).{8,20}$"))
            {
                MessageBox.Show("密码需要包含大小写字母、数字和特殊符号！");
                return;
            }
        }

        private void buttonCancel_Click(object sender, EventArgs e)
        {
            Close();
        }

        private void textBoxUserName_KeyPress(object sender, KeyPressEventArgs e)
        {
            if (!Regex.IsMatch(e.KeyChar.ToString(), @"^[0-9A-Za-z]$") && e.KeyChar != 8)
                e.Handled = true;
        }

        private void textBoxPassword_KeyPress(object sender, KeyPressEventArgs e)
        {
            if (!Regex.IsMatch(e.KeyChar.ToString(), @"^[0-9A-Za-z!@#$%^&*]$") && e.KeyChar != 8)
                e.Handled = true;
        }
    }
}
