import matplotlib.pyplot as plt
import seaborn as sb
import pandas as pd
import numpy as np

pre_intensity = pd.read_csv("initial_intensity.csv")
post_intensity = pd.read_csv("post_initial_intensity.csv")
coherence = pd.read_csv("degree_of_coherence.csv")

time = pre_intensity['time']
pre_I = pre_intensity['intensity']
post_I = post_intensity['intensity']
g_1 = coherence['coherence']

n = 1

sb.set()

plt.figure(n)
sb.lineplot(x=time, y=pre_I)
plt.xlabel("Time")
plt.ylabel("Intensity")
plt.savefig('pre_intensity.png')
n=n+1

plt.figure(n)
sb.lineplot(x=time, y=post_I)
plt.xlabel("Time")
plt.ylabel("Intensity")
plt.savefig('post_intensity.png')
n=n+1

plt.figure(n)
sb.lineplot(x=time, y=g_1)
plt.xlabel("Time")
plt.ylabel("Coherence")
plt.savefig('degree_of_coherence.png')
n=n+1

pre_intensity_fft = np.fft.fft(pre_intensity)
pre_magnitude_spectrum = np.abs(pre_intensity_fft)
pre_power_spectrum = np.abs(pre_intensity_fft) ** 2

post_intensity_fft = np.fft.fft(post_intensity)
post_magnitude_spectrum = np.abs(post_intensity_fft)
post_power_spectrum = np.abs(post_intensity_fft) ** 2


plt.figure(n)
sb.lineplot(pre_magnitude_spectrum)
plt.xlabel('Frequency')
plt.ylabel('Magnitude Spectrum')
plt.title('Fourier Transform of Intensity')
plt.savefig('pre_magnitude_spectrum.png')
n=n+1

plt.figure(n)
sb.lineplot(pre_power_spectrum)
plt.xlabel('Frequency')
plt.ylabel('Magnitude Spectrum')
plt.title('Fourier Transform of Intensity')
plt.savefig('pre_power_spectrum.png')
n=n+1

plt.figure(n)
sb.lineplot(post_magnitude_spectrum)
plt.xlabel('Frequency')
plt.ylabel('Magnitude Spectrum')
plt.title('Fourier Transform of Intensity')
plt.savefig('post_magnitude_spectrum.png')
n=n+1

plt.figure(n)
sb.lineplot(post_power_spectrum)
plt.xlabel('Frequency')
plt.ylabel('Magnitude Spectrum')
plt.title('Fourier Transform of Intensity')
plt.savefig('post_power_spectrum.png')
n=n+1