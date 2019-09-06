// starts resources to provision them.
build {

    aws_ami_from "aws-ubuntu-16.04" {
        name = "{{user `image_name`}}-aws-ubuntu-16.04"
        // this creates a new resource with settings inherited from the source  
    }

    aws_ami_from "vb-ubuntu-12.04" {
        name = "{{user `image_name`}}-vb-ubuntu-12.04"
        communicator "ssh" {
            ssh_username = "ubuntu"
        }
    }

    aws_ami_from "packer-vmw-ubuntu-16.04" {
        name = "{{user `image_name`}}-vmw-ubuntu-16.04"
    }

    provisioners {
        shell {
            inline = [
                "echo '{{user `my_secret`}}' :D"
            ]
        }

        shell {
            script = [
                "script-1.sh",
                "script-2.sh",
            ]
            override "vmware-iso" {
                execute_command = "echo 'password' | sudo -S bash {{.Path}}"
            }
        }

        upload "log.go" "/tmp" {
            timeout = "5s"
        }
    }

}

build {
    // build an ami using the ami from the previous build block.
    aws_ami_from "{{user `image_name`}}-aws-ubuntu-16.04" {
    }

    provisioners {
        shell {
            inline = [
                "HOLY GUACAMOLE !"
            ]
        }
    }
}