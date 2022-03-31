FROM archlinux
ADD . /home/go_project
RUN yes Y | pacman -Sy go && yes Y | pacman -Sy gcc && yes Y | pacman -Sy git
WORKDIR "/home/go_project"
ENV BITLY_OAUTH_TOKEN=448151789bf0264b0596dac054cdc900c10a1b40
CMD cd test && go test && echo "Done!"

