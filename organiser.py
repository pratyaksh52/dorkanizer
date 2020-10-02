# Program which will sort your downlaod folder on the basis of file types.

# importing libraries
# os for getting directories
# platform to get the OS being used by the user
# shutil to move files
# json to read from .json files
import os
import platform
import shutil
import json

# datetime module to add time stamp
# pytz for timezone
from datetime import datetime
import pytz


# checking if the user is not running windows
# terminating the application if not running windows
if platform.system() != "Windows":
    print("This program works only on windows, please use a windows machine or VM.")
    print("Platform: {}".format(platform.system()))
    exit()


# getting the path of the user's downloads folder
downloads_path = os.path.join(os.path.expanduser('~'), 'Downloads')

# function which returns the index of the beginning of the extension name
# iterates from the last index to reduce time
# returns None if it is a folder
def extension_index(filename):
    for i in range(-1, -len(filename), -1):
        if filename[i] == ".":
            return i + len(filename)
    
    return None



# Storing previous structure in a text file
# Contains date and time of creation for simplicity

# Creating the file or opening it if already exists
history = open(os.path.join(downloads_path, 'history.txt'), 'a+', encoding='utf-8')

# getting a list containing all file names in the directory
filenames = os.listdir(downloads_path)

# Adding the date and time stamp
date_and_time = datetime.now(pytz.timezone('Asia/Kolkata'))     # getting time and date
date = "Date: {}".format(date_and_time.strftime("%d/%m/%y"))    # storing date
time = "Time: {}".format(date_and_time.strftime("%H:%M:%S (%Z)"))   # storing time
history.write("\n    {}\n".format(date))    # writing the date
history.write("    {}\n\n".format(time))    # writing the time

# writing to the file
for name in filenames:  
    # removing 'history.txt' from the file list if exists
    if name == 'history.txt':
        filenames.remove('history.txt')
        continue
    
    # marking a folder if it is one
    if extension_index(name) == None:
        history.write("{} [*folder*]\n".format(name))   # showing it as folder
    else:
        history.write("{}\n".format(name))


history.write("\n\n{}\n".format("-"*100))   # adding a separator of 100 hyphens
history.close()     # closing the file



# reading and storing data containing all categories and extensions as dictionary
extensions_dict = json.load(open('extensions.json'))

# creating a dictionary of paths for the new directories
paths_dict = {}
for i in extensions_dict.keys():
    paths_dict[i] = os.path.join(downloads_path, i)


# Creating directories for categorization.
# If directory with same name already exists, new directory is not created.
# FileExistsError is raised when a directory with same name already exists.
for path in paths_dict.values():
    try:
        os.mkdir(path)      # creating the directory
    except FileExistsError:
        pass


# Moving the files
for name in filenames:
    # going to the next file if it is a folder
    if (extension_index(name) == None):
        continue

    extension = name[extension_index(name):]    # getting extenstion of the file
    source = os.path.join(downloads_path, name)     # storing the path of the file

    # loop to check for the extension category in the json and move the file accordingly
    for cat, ext_list in extensions_dict.items():
        # going to the next category if the current extension is not of this category
        if extension not in ext_list:
            continue
        else:
            destination = os.path.join(paths_dict[cat], name)   # setting destination of the file
            shutil.move(source, destination)    # moving the file to its category
            break   # breaking the loop to go to the next file
    else:
        # executed if the loop runs without breaking
        # if the loop runs withour breaking, means the file does not have a category
        destination = os.path.join(paths_dict['Miscellaneous'], name)   # setting destination to Miscellaneous
        shutil.move(source, destination)    # moving the file
